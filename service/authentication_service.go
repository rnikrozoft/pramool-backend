package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rnikrozoft/pramool-core/config"
	"github.com/rnikrozoft/pramool-core/model"
)

type AuthenticationService interface {
	GenerateToken(userID string) (string, error)
	LoginByTel(ctx context.Context, tel string) (string, error)
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

func (service authentication) GenerateToken(userID string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(service.appConfigs.Jwt.ExpireTime) * time.Hour)

	claims := &model.CustomClaims{
		UserID:   userID,
		LoggedIn: true,
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

func (service authentication) LoginByTel(ctx context.Context, tel string) (string, error) {
	sub, err := service.userService.FindUserIDByTelWithFallback(ctx, tel)
	if err != nil {
		return "", err
	}
	return service.GenerateToken(sub)
}
