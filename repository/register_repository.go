package repository

import (
	"context"

	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/uptrace/bun"
)

type Register interface {
	BeginTx(ctx context.Context) (bun.Tx, error)

	FindPhoneNumber(ctx context.Context, tel string) (*entity.TelVerify, error)
	RegisterTel(ctx context.Context, tel string) error
	RegisterTelIfNotExist(ctx context.Context, tel string) error
	RegisterUserWithTx(ctx context.Context, tx bun.Tx, user entity.User) error
}

type register struct {
	bun *bun.DB
}

func NewRegisterRepository(bun *bun.DB) Register {
	return register{
		bun: bun,
	}
}

func (r register) BeginTx(ctx context.Context) (bun.Tx, error) {
	return r.bun.BeginTx(ctx, nil)
}

func (r register) FindPhoneNumber(ctx context.Context, tel string) (*entity.TelVerify, error) {
	table := new(entity.TelVerify)
	query := `SELECT tel FROM tel_verify WHERE tel = ?`
	err := r.bun.NewRaw(query, tel).Scan(ctx, table)
	if err != nil {
		return nil, err
	}
	return table, nil
}

func (r register) RegisterTel(ctx context.Context, tel string) error {
	query := `INSERT INTO tel_verify (tel) VALUES (?)`
	_, err := r.bun.NewRaw(query, tel).Exec(ctx)
	return err
}

func (r register) RegisterTelIfNotExist(ctx context.Context, tel string) error {
	_, err := r.bun.NewRaw(`INSERT INTO tel_verify (tel) VALUES (?) ON CONFLICT (tel) DO NOTHING`, tel).Exec(ctx)
	return err
}

func (r register) RegisterUserWithTx(ctx context.Context, tx bun.Tx, user entity.User) error {
	query := `
	INSERT INTO users (
		user_id, tel, email, facebook, bank_id, bank_account_name, bank_account_number,
		first_name, last_name, address_primary, address, soi, road, sub_district, district, province, zip_code
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := tx.NewRaw(query,
		user.UserID,
		user.Tel,
		user.Email,
		user.Facebook,
		user.BankID,
		user.BankAccountName,
		user.BankAccountNumber,
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
