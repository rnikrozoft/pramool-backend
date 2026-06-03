package mapping

import (
	"strings"

	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/model/entity"
)

func ToUserEntity(req dto.UserRegisterRequest) entity.User {
	nationalID := strings.TrimSpace(req.NationalID)
	if nationalID == "" {
		nationalID = strings.TrimSpace(req.UserID)
	}
	return entity.User{
		NationalID:        nationalID,
		Tel:               req.Tel,
		Email:             req.Email,
		Facebook:          req.Facebook,
		BankID:            req.BankID,
		BankAccountName:   req.BankAccountName,
		BankAccountNumber: req.BankAccountNumber,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		AddressPrimary:    req.AddressPrimary,
		Address:           req.Address,
		Soi:               req.Soi,
		Road:              req.Road,
		SubDistrict:       req.SubDistrict,
		District:          req.District,
		Province:          req.Province,
		ZipCode:           req.ZipCode,
		Credit:            req.Credit,
	}
}

func ToUserDTO(entity entity.User) dto.UserRegisterRequest {
	return dto.UserRegisterRequest{
		NationalID:        entity.NationalID,
		Tel:               entity.Tel,
		Email:             entity.Email,
		Facebook:          entity.Facebook,
		BankID:            entity.BankID,
		BankAccountName:   entity.BankAccountName,
		BankAccountNumber: entity.BankAccountNumber,
		FirstName:         entity.FirstName,
		LastName:          entity.LastName,
		AddressPrimary:    entity.AddressPrimary,
		Address:           entity.Address,
		Soi:               entity.Soi,
		Road:              entity.Road,
		SubDistrict:       entity.SubDistrict,
		District:          entity.District,
		Province:          entity.Province,
		ZipCode:           entity.ZipCode,
		Credit:            entity.Credit,
	}
}
