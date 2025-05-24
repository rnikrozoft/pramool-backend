package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/rnikrozoft/pramool.in.th-backend/mapping"
	"github.com/rnikrozoft/pramool.in.th-backend/model/dto"
	"github.com/rnikrozoft/pramool.in.th-backend/repository"
)

type RegisterService interface {
	Register(ctx context.Context, user dto.User) (string, error)
}

type register struct {
	registerRepository    repository.Register
	authenticationService AuthenticationService
}

func NewRegisterService(
	registerRepository repository.Register,
	authenticationService AuthenticationService,
) RegisterService {
	return register{
		registerRepository:    registerRepository,
		authenticationService: authenticationService,
	}
}

func (service register) Register(ctx context.Context, user dto.User) (string, error) {
	r := mapping.ToUserEntity(user)

	h := sha256.Sum256([]byte(user.Password))
	r.Password = hex.EncodeToString(h[:])

	if err := service.registerRepository.Register(ctx, r); err != nil {
		return "", err
	}

	token, err := service.authenticationService.Login(ctx, r.Email, r.Password)
	if err != nil {
		return "", err
	}
	return token, nil
}
