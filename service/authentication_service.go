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
	ValidateRefreshToken(tokenString string) (userID string, err error)
	LoginByTel(ctx context.Context, tel, password string) (LoginTokens, error)
	Login(ctx context.Context, identifier, password string) (LoginTokens, error)
}

type authentication struct {
	appConfigs  config.AppConfigs
	userService UserService
}

func NewAuthenticationService(appConfigs config.AppConfigs, userService UserService) AuthenticationService {
	return authentication{
		appConfigs:  appConfigs,
		userService: userService,
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
	return service.sign(userID, model.TokenUseAccess, service.accessTTL())
}

func (service authentication) GenerateRefreshToken(userID string) (string, error) {
	return service.sign(userID, model.TokenUseRefresh, service.refreshTTL())
}

func (service authentication) sign(userID, use string, ttl time.Duration) (string, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return "", errors.New("empty user id")
	}
	expirationTime := time.Now().Add(ttl)
	claims := &model.CustomClaims{
		UserID:   userID,
		LoggedIn: true,
		TokenUse: use,
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

func (service authentication) ValidateRefreshToken(tokenString string) (string, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return "", errNotRefreshToken
	}
	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.appConfigs.Jwt.Secret), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	claims, ok := token.Claims.(*model.CustomClaims)
	if !ok {
		return "", errNotRefreshToken
	}
	if claims.TokenUse != model.TokenUseRefresh {
		return "", errNotRefreshToken
	}
	userID := strings.TrimSpace(claims.UserID)
	if userID == "" {
		userID = strings.TrimSpace(claims.Subject)
	}
	if userID == "" {
		return "", errNotRefreshToken
	}
	return userID, nil
}

func (service authentication) Login(ctx context.Context, identifier, password string) (LoginTokens, error) {
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
	return service.LoginByTel(ctx, tel, password)
}

func (service authentication) LoginByTel(ctx context.Context, tel, password string) (LoginTokens, error) {
	var empty LoginTokens
	tel = strings.TrimSpace(tel)
	password = strings.TrimSpace(password)
	hash, err := service.userService.GetLoginPasswordHash(ctx, tel)
	if err != nil {
		return empty, err
	}
	if strings.TrimSpace(hash) != "" {
		if password == "" {
			return empty, errPasswordRequired()
		}
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
			return empty, errInvalidCredentials()
		}
	}
	sub, err := service.userService.FindUserIDByTelWithFallback(ctx, tel)
	if err != nil {
		return empty, err
	}
	access, err := service.GenerateAccessToken(sub)
	if err != nil {
		return empty, err
	}
	refresh, err := service.GenerateRefreshToken(sub)
	if err != nil {
		return empty, err
	}
	return LoginTokens{Access: access, Refresh: refresh}, nil
}
