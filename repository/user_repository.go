package repository

import (
	"context"

	"github.com/rnikrozoft/pramool.in.th-backend/model/entity"
	"github.com/uptrace/bun"
)

type UserRepository interface {
	FindByTel(ctx context.Context, tel string) (*entity.User, error)
	GetMyInformation(ctx context.Context, userID string) (*entity.User, error)
}

type user struct {
	bun *bun.DB
}

func NewUserRepository(bun *bun.DB) UserRepository {
	return user{
		bun: bun,
	}
}

func (r user) FindByTel(ctx context.Context, tel string) (*entity.User, error) {
	user := new(entity.User)
	query := `SELECT tel FROM users WHERE tel = ?`
	err := r.bun.NewRaw(query, tel).Scan(ctx, user)
	return user, err
}

func (r user) GetMyInformation(ctx context.Context, userID string) (*entity.User, error) {
	user := new(entity.User)
	query := `SELECT * FROM users WHERE user_id = ?`
	err := r.bun.NewRaw(query, userID).Scan(ctx, user)
	return user, err
}
