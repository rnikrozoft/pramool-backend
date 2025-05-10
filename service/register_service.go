package service

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/rnikrozoft/pramool.in.th-backend/model"
	"github.com/rnikrozoft/pramool.in.th-backend/model/entity"
	"github.com/rnikrozoft/pramool.in.th-backend/repository"
)

type Register interface {
	Register(userData model.User) error
}

type register struct {
	registerRepository repository.Register
}

func NewRegisterService(registerRepositryo repository.Register) Register {
	return register{
		registerRepository: registerRepositryo,
	}
}

func (s register) Register(userData model.User) error {
	r := entity.User{
		Email: userData.Email,
	}

	h := sha256.Sum256([]byte(userData.Password))
	r.Password = hex.EncodeToString(h[:])

	if err := s.registerRepository.Register(r); err != nil {
		return err
	}
	return nil
}
