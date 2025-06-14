package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rnikrozoft/pramool.in.th-backend/exception"
	"github.com/rnikrozoft/pramool.in.th-backend/model/entity"
	"github.com/rnikrozoft/pramool.in.th-backend/repository"
)

type UserService interface {
	IsTelAlreadyUsed(ctx context.Context, tel string) (bool, error)
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

func (service user) IsTelAlreadyUsed(ctx context.Context, tel string) (bool, error) {
	data, err := service.userRepository.FindByTel(ctx, tel)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ไม่เจอเบอร์โทรใน DB = ยังไม่ถูกใช้
			return false, nil
		}
		return false, exception.Internal(err)
	}
	// เจอข้อมูล = เบอร์โทรถูกใช้แล้ว
	return data != nil, nil
}

func (service user) GetMyInfo(ctx context.Context, userID string) (*entity.User, error) {
	myInfo, err := service.userRepository.GetMyInformation(ctx, userID)
	if err != nil {
		return nil, err
	}
	return myInfo, nil
}
