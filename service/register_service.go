package service

import (
	"context"

	"github.com/rnikrozoft/pramool.in.th-backend/model/entity"
	"github.com/rnikrozoft/pramool.in.th-backend/repository"
)

type RegisterService interface {
	Register(ctx context.Context, user entity.User) error
}

type register struct {
	registerRepository repository.Register
}

func NewRegisterService(
	registerRepository repository.Register,
) RegisterService {
	return register{
		registerRepository: registerRepository,
	}
}

func (service register) Register(ctx context.Context, user entity.User) error {
	if err := service.registerRepository.RegisterUser(ctx, user); err != nil {
		return err
	}
	return nil
}
