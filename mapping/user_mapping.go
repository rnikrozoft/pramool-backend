package mapping

import (
	"github.com/rnikrozoft/pramool.in.th-backend/model/dto"
	"github.com/rnikrozoft/pramool.in.th-backend/model/entity"
)

func ToUserEntity(req dto.UserRegisterRequest) entity.User {
	return entity.User{
		UserID:         req.UserID,
		Tel:            req.Tel,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		AddressPrimary: req.AddressPrimary,
		Address:        req.Address,
		Soi:            req.Soi,
		Road:           req.Road,
		SubDistrict:    req.SubDistrict,
		District:       req.District,
		Province:       req.Province,
		ZipCode:        req.ZipCode,
	}
}

func ToUserDTO(entity entity.User) dto.UserRegisterRequest {
	return dto.UserRegisterRequest{
		UserID:         entity.UserID,
		Tel:            entity.Tel,
		FirstName:      entity.FirstName,
		LastName:       entity.LastName,
		AddressPrimary: entity.AddressPrimary,
		Address:        entity.Address,
		Soi:            entity.Soi,
		Road:           entity.Road,
		SubDistrict:    entity.SubDistrict,
		District:       entity.District,
		Province:       entity.Province,
		ZipCode:        entity.ZipCode,
	}
}
