package repository

import (
	"github.com/rnikrozoft/pramool.in.th-backend/model/entity"
	"gorm.io/gorm"
)

type Register interface {
	Register(userData entity.User) error
}

type register struct {
	gorm *gorm.DB
}

func NewRegisterRepository(gorm *gorm.DB) Register {
	return register{
		gorm: gorm,
	}
}

func (r register) Register(userData entity.User) error {
	result := r.gorm.Create(&userData)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
