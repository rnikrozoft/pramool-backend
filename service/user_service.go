package service

import (
	"context"

	"github.com/rnikrozoft/pramool.in.th-backend/model/entity"
	"github.com/rnikrozoft/pramool.in.th-backend/repository"
)

type UserService interface {
	GetMyInfo(ctx context.Context, userID string) (*entity.User, error)
}

type user struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return user{
		userRepository: userRepository,
	}
}

func (service user) GetMyInfo(ctx context.Context, userID string) (*entity.User, error) {
	myInfo, err := service.userRepository.GetMyInformation(ctx, userID)
	if err != nil {
		return nil, err
	}
	return myInfo, nil
}
