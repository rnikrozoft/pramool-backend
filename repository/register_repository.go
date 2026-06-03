package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/rnikrozoft/pramool-core/internal/nationalid"
	"github.com/uptrace/bun"
)

type Register interface {
	BeginTx(ctx context.Context) (bun.Tx, error)

	FindPhoneNumber(ctx context.Context, tel string) (*entity.TelVerify, error)
	RegisterTel(ctx context.Context, tel string) error
	RegisterTelIfNotExist(ctx context.Context, tel string) error
	RegisterUserWithTx(ctx context.Context, tx bun.Tx, user entity.User) error

	UpsertTelVerifySignup(ctx context.Context, firstName, lastName, tel, email, passwordHash string) error
	GetPasswordHashFromTelVerify(ctx context.Context, tel string) (string, error)
	// TelVerifyHasPassword is true when this tel already completed signup step (password set); must login, not POST /auth/signup again.
	TelVerifyHasPassword(ctx context.Context, tel string) (bool, error)
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
		national_id_hash, national_id_enc, tel, email, facebook, bank_id, bank_account_name, bank_account_number,
		first_name, last_name, address_primary, address, soi, road, sub_district, district, province, zip_code,
		password_hash
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	hash, enc, err := nationalid.PrepareStorage(user.NationalID)
	if err != nil {
		return err
	}
	_, err = tx.NewRaw(query,
		nullIfEmpty(hash),
		nullIfEmpty(enc),
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
		nullIfEmpty(user.PasswordHash),
	).Exec(ctx)
	return err
}

func nullIfEmpty(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return strings.TrimSpace(s)
}

func (r register) UpsertTelVerifySignup(ctx context.Context, firstName, lastName, tel, email, passwordHash string) error {
	_, err := r.bun.NewRaw(`
		INSERT INTO tel_verify (tel, signup_first_name, signup_last_name, signup_email, password_hash)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (tel) DO UPDATE SET
			signup_first_name = EXCLUDED.signup_first_name,
			signup_last_name = EXCLUDED.signup_last_name,
			signup_email = EXCLUDED.signup_email,
			password_hash = EXCLUDED.password_hash,
			updated_at = NOW()
	`, tel, firstName, lastName, email, passwordHash).Exec(ctx)
	return err
}

func (r register) GetPasswordHashFromTelVerify(ctx context.Context, tel string) (string, error) {
	var h sql.NullString
	err := r.bun.NewRaw(`SELECT password_hash FROM tel_verify WHERE tel = ?`, tel).Scan(ctx, &h)
	if err != nil {
		return "", err
	}
	if !h.Valid {
		return "", nil
	}
	return h.String, nil
}

func (r register) TelVerifyHasPassword(ctx context.Context, tel string) (bool, error) {
	tel = strings.TrimSpace(tel)
	if tel == "" {
		return false, nil
	}
	var h sql.NullString
	err := r.bun.NewRaw(`SELECT password_hash FROM tel_verify WHERE TRIM(tel) = ?`, tel).Scan(ctx, &h)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if !h.Valid {
		return false, nil
	}
	return strings.TrimSpace(h.String) != "", nil
}
