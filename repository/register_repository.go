package repository

import (
	"context"

	"github.com/rnikrozoft/pramool.in.th-backend/model/entity"
	"github.com/uptrace/bun"
)

type Register interface {
	RegisterUser(ctx context.Context, user entity.User) error
}

type register struct {
	bun *bun.DB
}

func NewRegisterRepository(bun *bun.DB) Register {
	return register{
		bun: bun,
	}
}

func (r register) RegisterUser(ctx context.Context, user entity.User) error {
	query := `
	INSERT INTO users (
		user_id, tel, first_name, last_name, address_primary,
		address, soi, road, sub_district, district, province, zip_code
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.bun.NewRaw(query,
		user.UserID,
		user.Tel,
		user.FirstName,
		user.LastName,
		user.AddressPrimary,
		user.Address,
		user.Soi,
		user.Road,
		user.SubDistrict,
		user.District,
		user.Province,
		user.ZipCode,
	).Exec(ctx)
	return err
}
