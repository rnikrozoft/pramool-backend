package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rnikrozoft/pramool.in.th-backend/config"
	"github.com/rnikrozoft/pramool.in.th-backend/model"
	"github.com/rnikrozoft/pramool.in.th-backend/repository"
)

type AuthenticationService interface {
	GenerateToken(userID string) (string, error)
	LoginByTel(ctx context.Context, tel string) (string, error)
}

type authentication struct {
	appConfigs     config.AppConfigs
	userRepository repository.UserRepository
}

func NewAuthenticationService(appConfigs config.AppConfigs, userRepository repository.UserRepository) AuthenticationService {
	return authentication{
		appConfigs:     appConfigs,
		userRepository: userRepository,
	}
}

func (service authentication) GenerateToken(userID string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(service.appConfigs.Jwt.ExpireTime) * time.Hour)

	claims := &model.CustomClaims{
		UserID:   userID,
		LoggedIn: true,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    service.appConfigs.Jwt.Issuer,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(service.appConfigs.Jwt.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (service authentication) LoginByTel(ctx context.Context, tel string) (string, error) {
	userData, err := service.userRepository.FindByTel(ctx, tel)
	if err != nil {
		return "", err
	}

	token, err := service.GenerateToken(userData.UserID)
	if err != nil {
		return "", err
	}
	return token, nil
}
