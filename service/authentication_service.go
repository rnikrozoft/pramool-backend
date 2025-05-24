package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rnikrozoft/pramool.in.th-backend/config"
	"github.com/rnikrozoft/pramool.in.th-backend/repository"
)

type AuthenticationService interface {
	Login(ctx context.Context, email, password string) (string, error)
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

func (service authentication) Login(ctx context.Context, email, password string) (string, error) {
	userId, err := service.userRepository.FindUserIdByEmailAndPassword(ctx, email, password)
	if err != nil {
		return "", err
	}

	if userId == "" {
		return "", errors.New("user not found")
	}

	token, err := service.generateToken(userId)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (service authentication) generateToken(userID string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(service.appConfigs.Jwt.ExpireTime) * time.Hour)

	type CustomClaims struct {
		UserID string `json:"user_id"`
		jwt.RegisteredClaims
	}

	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
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
