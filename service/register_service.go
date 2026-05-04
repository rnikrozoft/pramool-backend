package service

import (
	"context"
	"strings"

	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/rnikrozoft/pramool-core/repository"
	"golang.org/x/crypto/bcrypt"
)

type RegisterService interface {
	RegisterTelIfNotExist(ctx context.Context, tel string) error
	RegisterUser(ctx context.Context, user entity.User) error
	SignupWithPassword(ctx context.Context, firstName, lastName, tel, email, plainPassword string) error
	TelVerifyHasPassword(ctx context.Context, tel string) (bool, error)
}

type register struct {
	registerRepository repository.Register
	userService        UserService
}

func NewRegisterService(
	registerRepository repository.Register,
	userService UserService,
) RegisterService {
	return register{
		registerRepository: registerRepository,
		userService:        userService,
	}
}

func (service register) RegisterTelIfNotExist(ctx context.Context, tel string) error {
	return service.registerRepository.RegisterTelIfNotExist(ctx, tel)
}

func (service register) TelVerifyHasPassword(ctx context.Context, tel string) (bool, error) {
	return service.registerRepository.TelVerifyHasPassword(ctx, tel)
}

func (service register) RegisterUser(ctx context.Context, user entity.User) error {
	user.UserID = strings.TrimSpace(user.UserID)
	user.Tel = strings.TrimSpace(user.Tel)
	user.Email = strings.TrimSpace(user.Email)

	idTaken, err := service.userService.ExistsRegisteredUserID(ctx, user.UserID)
	if err != nil {
		return err
	}
	if idTaken {
		return ErrNationalIDAlreadyRegistered
	}

	telUsed, err := service.userService.IsTelAlreadyUsed(ctx, user.Tel)
	if err != nil {
		return err
	}
	if telUsed {
		return ErrTelHasFullUserRecord
	}

	if user.Email != "" {
		emailTaken, err := service.userService.IsEmailTakenByOtherTel(ctx, user.Email, user.Tel)
		if err != nil {
			return err
		}
		if emailTaken {
			return ErrEmailAlreadyRegistered
		}
	}

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

func (service register) SignupWithPassword(ctx context.Context, firstName, lastName, tel, email, plainPassword string) error {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	tel = strings.TrimSpace(tel)
	email = strings.TrimSpace(email)
	if email != "" {
		taken, err := service.userService.IsEmailTakenByOtherTel(ctx, email, tel)
		if err != nil {
			return err
		}
		if taken {
			return ErrEmailAlreadyRegistered
		}
	}
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return service.registerRepository.UpsertTelVerifySignup(ctx, firstName, lastName, tel, email, string(hashBytes))
}
