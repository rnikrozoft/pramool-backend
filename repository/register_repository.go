package repository

import (
	"context"

	"github.com/rnikrozoft/pramool.in.th-backend/model/entity"
	"github.com/uptrace/bun"
)

type Register interface {
	Register(ctx context.Context, user entity.User) error
}

type register struct {
	bun *bun.DB
}

func NewRegisterRepository(bun *bun.DB) Register {
	return register{
		bun: bun,
	}
}

func (r register) Register(ctx context.Context, user entity.User) error {
	query := `
		INSERT INTO users (
			user_id, 
			email, 
			password,
			first_name,
			last_name
		) VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.bun.NewRaw(query, user.UserId, user.Email, user.Password, user.FirstName, user.LastName).Exec(ctx)
	return err
}
