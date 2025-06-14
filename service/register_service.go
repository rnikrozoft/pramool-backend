package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rnikrozoft/pramool.in.th-backend/model/entity"
	"github.com/rnikrozoft/pramool.in.th-backend/repository"
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
	data, err := service.registerRepository.FindPhoneNumber(ctx, tel)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if data == nil || errors.Is(err, sql.ErrNoRows) {
		if err := service.registerRepository.RegisterTel(ctx, tel); err != nil {
			return err
		}
	}
	return nil
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

	if err := service.registerRepository.SetTelIsVerifyWithTx(ctx, tx, user.Tel); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
