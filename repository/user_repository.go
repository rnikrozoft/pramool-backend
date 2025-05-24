package repository

import (
	"context"

	"github.com/rnikrozoft/pramool.in.th-backend/model/entity"
	"github.com/uptrace/bun"
)

type UserRepository interface {
	FindUserIdByEmailAndPassword(ctx context.Context, email, password string) (string, error)
}

type user struct {
	bun *bun.DB
}

func NewUserRepository(bun *bun.DB) UserRepository {
	return user{
		bun: bun,
	}
}

func (r user) FindUserIdByEmailAndPassword(ctx context.Context, email, password string) (string, error) {
	user := new(entity.User)
	query := `SELECT user_id FROM users WHERE email = ? AND password = ?`
	err := r.bun.NewRaw(query, email, password).Scan(ctx, user)
	return user.UserId, err
}
