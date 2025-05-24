package mapping

import (
	"github.com/rnikrozoft/pramool.in.th-backend/model/dto"
	"github.com/rnikrozoft/pramool.in.th-backend/model/entity"
)

func ToUserEntity(dto dto.User) entity.User {
	return entity.User{
		UserId:    dto.UserId,
		Email:     dto.Email,
		Password:  dto.Password,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
	}
}

func ToUserDTO(entity entity.User) dto.User {
	return dto.User{
		UserId:    entity.UserId,
		Email:     entity.Email,
		Password:  entity.Password,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
	}
}
