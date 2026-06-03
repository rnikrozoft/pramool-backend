package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/rnikrozoft/pramool-core/config"
	"github.com/rnikrozoft/pramool-core/exception"
	"github.com/rnikrozoft/pramool-core/model"
)

// LoginTokens holds access + refresh JWTs returned after login.
type LoginTokens struct {
	Access  string
	Refresh string
}

var errNotRefreshToken = errors.New("not a refresh token")

type AuthenticationService interface {
	// GenerateToken issues an access token (alias for GenerateAccessToken).
	GenerateToken(userID string) (string, error)
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	GenerateRefreshTokenWithRemember(userID string, remember bool) (string, error)
	ValidateRefreshToken(tokenString string) (userID string, err error)
	ParseRefreshToken(tokenString string) (userID string, remember bool, err error)
	LoginByTel(ctx context.Context, tel, password string, remember bool) (LoginTokens, error)
	Login(ctx context.Context, identifier, password string, remember bool) (LoginTokens, error)
}

type authentication struct {
	appConfigs     config.AppConfigs
	userService    UserService
	privacyService PrivacyService
}

func NewAuthenticationService(appConfigs config.AppConfigs, userService UserService, privacyService PrivacyService) AuthenticationService {
	return authentication{
		appConfigs:     appConfigs,
		userService:    userService,
		privacyService: privacyService,
	}
}

func (service authentication) accessTTL() time.Duration {
	h := service.appConfigs.Jwt.ExpireTime
	if h <= 0 {
		h = 1
	}
	return time.Duration(h) * time.Hour
}

func (service authentication) refreshTTL() time.Duration {
	h := service.appConfigs.Jwt.RefreshExpireTime
	if h <= 0 {
		h = 168 // 7 days
	}
	return time.Duration(h) * time.Hour
}

func (service authentication) GenerateToken(userID string) (string, error) {
	return service.GenerateAccessToken(userID)
}

func (service authentication) GenerateAccessToken(userID string) (string, error) {
	return service.sign(userID, model.TokenUseAccess, service.accessTTL(), false)
}

func (service authentication) GenerateRefreshToken(userID string) (string, error) {
	return service.GenerateRefreshTokenWithRemember(userID, true)
}

func (service authentication) GenerateRefreshTokenWithRemember(userID string, remember bool) (string, error) {
	ttl := service.accessTTL()
	if remember {
		ttl = service.refreshTTL()
	}
	return service.sign(userID, model.TokenUseRefresh, ttl, remember)
}

func (service authentication) sign(userID, use string, ttl time.Duration, remember bool) (string, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return "", errors.New("empty user id")
	}
	expirationTime := time.Now().Add(ttl)
	claims := &model.CustomClaims{
		UserID:     userID,
		LoggedIn:   true,
		TokenUse:   use,
		RememberMe: remember && use == model.TokenUseRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    service.appConfigs.Jwt.Issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(service.appConfigs.Jwt.Secret))
}

func (service authentication) ParseRefreshToken(tokenString string) (string, bool, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return "", false, errNotRefreshToken
	}
	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.appConfigs.Jwt.Secret), nil
	})
	if err != nil || !token.Valid {
		return "", false, err
	}
	claims, ok := token.Claims.(*model.CustomClaims)
	if !ok {
		return "", false, errNotRefreshToken
	}
	if claims.TokenUse != model.TokenUseRefresh {
		return "", false, errNotRefreshToken
	}
	userID := strings.TrimSpace(claims.UserID)
	if userID == "" {
		userID = strings.TrimSpace(claims.Subject)
	}
	if userID == "" {
		return "", false, errNotRefreshToken
	}
	remember := claims.RememberMe
	if !remember && claims.ExpiresAt != nil && claims.IssuedAt != nil {
		// Legacy refresh tokens (before remember_me claim): infer from lifetime.
		if claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time) > 2*service.accessTTL() {
			remember = true
		}
	}
	return userID, remember, nil
}

func (service authentication) ValidateRefreshToken(tokenString string) (string, error) {
	userID, _, err := service.ParseRefreshToken(tokenString)
	return userID, err
}

func (service authentication) Login(ctx context.Context, identifier, password string, remember bool) (LoginTokens, error) {
	var empty LoginTokens
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return empty, exception.BadRequest(errors.New("กรุณากรอกเบอร์โทรศัพท์หรืออีเมล"))
	}
	tel, err := service.userService.ResolveLoginIdentifier(ctx, identifier)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return empty, errInvalidCredentials()
		}
		return empty, err
	}
	return service.LoginByTel(ctx, tel, password, remember)
}

func (service authentication) LoginByTel(ctx context.Context, tel, password string, remember bool) (LoginTokens, error) {
	var empty LoginTokens
	tel = strings.TrimSpace(tel)
	password = strings.TrimSpace(password)
	hash, err := service.userService.GetLoginPasswordHash(ctx, tel)
	if err != nil {
		return empty, err
	}
	if strings.TrimSpace(hash) == "" {
		return empty, errInvalidCredentials()
	}
	if password == "" {
		return empty, errPasswordRequired()
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return empty, errInvalidCredentials()
	}
	sub, err := service.userService.FindUserIDByTelWithFallback(ctx, tel)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return empty, errInvalidCredentials()
		}
		return empty, err
	}
	if service.privacyService != nil {
		deleted, derr := service.privacyService.IsAccountDeletedByTel(ctx, tel)
		if derr == nil && deleted {
			return empty, exception.BadRequest(errors.New("บัญชีนี้ถูกลบแล้ว"))
		}
	}
	access, err := service.GenerateAccessToken(sub)
	if err != nil {
		return empty, err
	}
	refresh, err := service.GenerateRefreshTokenWithRemember(sub, remember)
	if err != nil {
		return empty, err
	}
	return LoginTokens{Access: access, Refresh: refresh}, nil
}
