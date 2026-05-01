package service

import (
	"context"

	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/rnikrozoft/pramool-core/repository"
)

type RegisterService interface {
	RegisterTelIfNotExist(ctx context.Context, tel string) error
	RegisterUser(ctx context.Context, user entity.User) error
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

func (service register) RegisterTelIfNotExist(ctx context.Context, tel string) error {
	return service.registerRepository.RegisterTelIfNotExist(ctx, tel)
}

func (service register) RegisterUser(ctx context.Context, user entity.User) error {
	tx, err := service.registerRepository.BeginTx(ctx)
	if err != nil {
		return err
	}

	if err := service.registerRepository.RegisterUserWithTx(ctx, tx, user); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
