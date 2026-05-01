package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/uptrace/bun"
)

// ErrNoUserUpdated is returned when UPDATE users matched no row (e.g. wrong user_id).
var ErrNoUserUpdated = errors.New("no user row updated")

type UserRepository interface {
	FindByTel(ctx context.Context, tel string) (*entity.User, error)
	FindTelVerifyByTel(ctx context.Context, tel string) (*entity.TelVerify, error)
	GetMyInformation(ctx context.Context, userID string) (*entity.User, error)
	AddCreditByUserID(ctx context.Context, userID string, amount int64) error

	FindUserIDByTel(ctx context.Context, tel string) (string, error)
	FindByID(ctx context.Context, userID string) (*entity.User, error)
	Upsert(ctx context.Context, u *entity.User) error
	IsTelAlreadyUsed(ctx context.Context, tel string) (bool, error)
	IsTelUsedByOtherUser(ctx context.Context, userID, tel string) (bool, error)
	HasUserRecordForSubject(ctx context.Context, subject string) (bool, error)
	ResetExpiredOTPBanIfNeeded(ctx context.Context, tel string) error
	SelectOTPBanUntil(ctx context.Context, tel string) (*time.Time, error)
	RecordOTPTimeout(ctx context.Context, tel string) (*time.Time, int, error)
	UpdateProfile(ctx context.Context, p entity.ProfileUpdate) error
}

type user struct {
	bun *bun.DB
}

func NewUserRepository(bun *bun.DB) UserRepository {
	return user{bun: bun}
}

func (r user) FindByTel(ctx context.Context, tel string) (*entity.User, error) {
	u := new(entity.User)
	query := `SELECT user_id, tel FROM users WHERE tel = ?`
	err := r.bun.NewRaw(query, tel).Scan(ctx, u)
	return u, err
}

func (r user) FindTelVerifyByTel(ctx context.Context, tel string) (*entity.TelVerify, error) {
	tv := new(entity.TelVerify)
	query := `SELECT tel FROM tel_verify WHERE tel = ?`
	err := r.bun.NewRaw(query, tel).Scan(ctx, tv)
	return tv, err
}

func (r user) GetMyInformation(ctx context.Context, userID string) (*entity.User, error) {
	return r.FindByID(ctx, userID)
}

func (r user) AddCreditByUserID(ctx context.Context, userID string, amount int64) error {
	query := `UPDATE users SET credit = credit + ?, updated_at = NOW() WHERE user_id = ?`
	_, err := r.bun.NewRaw(query, amount, userID).Exec(ctx)
	return err
}

func (r user) FindUserIDByTel(ctx context.Context, tel string) (string, error) {
	var id string
	err := r.bun.NewRaw(`SELECT user_id FROM users WHERE tel = ?`, tel).Scan(ctx, &id)
	return id, err
}

func (r user) FindByID(ctx context.Context, userID string) (*entity.User, error) {
	u := new(entity.User)
	err := r.bun.NewRaw(`
		SELECT user_id, tel,
		       COALESCE(first_name, '') AS first_name,
		       COALESCE(last_name, '') AS last_name,
		       COALESCE(address_primary, '') AS address_primary,
		       COALESCE(address, '') AS address,
		       COALESCE(soi, '') AS soi,
		       COALESCE(road, '') AS road,
		       COALESCE(sub_district, '') AS sub_district,
		       COALESCE(district, '') AS district,
		       COALESCE(province, '') AS province,
		       COALESCE(zip_code, '') AS zip_code,
		       COALESCE(email, '') AS email,
		       COALESCE(facebook, '') AS facebook,
		       COALESCE(credit, 0) AS credit,
		       created_at, updated_at
		FROM users WHERE user_id = ?
	`, userID).Scan(ctx, u)
	return u, err
}

func (r user) Upsert(ctx context.Context, u *entity.User) error {
	_, err := r.bun.NewRaw(`
		INSERT INTO users (
			user_id, tel, email, facebook, first_name, last_name,
			address_primary, address, soi, road, sub_district, district, province, zip_code
		)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT (user_id) DO UPDATE SET
			email = EXCLUDED.email,
			facebook = EXCLUDED.facebook,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			address_primary = EXCLUDED.address_primary,
			address = EXCLUDED.address,
			soi = EXCLUDED.soi,
			road = EXCLUDED.road,
			sub_district = EXCLUDED.sub_district,
			district = EXCLUDED.district,
			province = EXCLUDED.province,
			zip_code = EXCLUDED.zip_code,
			updated_at = NOW()
	`, u.UserID, u.Tel, u.Email, u.Facebook, u.FirstName, u.LastName,
		u.AddressPrimary, u.Address, u.Soi, u.Road, u.SubDistrict, u.District, u.Province, u.ZipCode).Exec(ctx)
	return err
}

func (r user) IsTelAlreadyUsed(ctx context.Context, tel string) (bool, error) {
	var exists bool
	err := r.bun.NewRaw(`SELECT EXISTS(SELECT 1 FROM users WHERE tel = ?)`, tel).Scan(ctx, &exists)
	return exists, err
}

func (r user) IsTelUsedByOtherUser(ctx context.Context, userID, tel string) (bool, error) {
	var exists bool
	err := r.bun.NewRaw(`SELECT EXISTS(SELECT 1 FROM users WHERE tel = ? AND user_id <> ?)`, tel, userID).Scan(ctx, &exists)
	return exists, err
}

func (r user) HasUserRecordForSubject(ctx context.Context, subject string) (bool, error) {
	var exists bool
	err := r.bun.NewRaw(`
		SELECT EXISTS(
			SELECT 1 FROM users WHERE user_id = ? OR tel = ?
		)
	`, subject, subject).Scan(ctx, &exists)
	return exists, err
}

func (r user) ResetExpiredOTPBanIfNeeded(ctx context.Context, tel string) error {
	_, err := r.bun.NewRaw(`
		UPDATE tel_verify
		SET otp_timeout_count = 0, otp_banned_until = NULL
		WHERE tel = ? AND otp_banned_until IS NOT NULL AND otp_banned_until <= NOW()
	`, tel).Exec(ctx)
	return err
}

func (r user) SelectOTPBanUntil(ctx context.Context, tel string) (*time.Time, error) {
	var bannedUntil sql.NullTime
	err := r.bun.NewRaw(`SELECT otp_banned_until FROM tel_verify WHERE tel = ?`, tel).Scan(ctx, &bannedUntil)
	if err != nil {
		return nil, err
	}
	if !bannedUntil.Valid {
		return nil, nil
	}
	t := bannedUntil.Time
	return &t, nil
}

func (r user) RecordOTPTimeout(ctx context.Context, tel string) (*time.Time, int, error) {
	var timeoutCount int
	var bannedUntil sql.NullTime
	err := r.bun.NewRaw(`
		UPDATE tel_verify
		SET otp_timeout_count = CASE
				WHEN otp_banned_until IS NULL OR otp_banned_until <= NOW() THEN otp_timeout_count + 1
				ELSE otp_timeout_count
			END,
			otp_banned_until = CASE
				WHEN (otp_banned_until IS NULL OR otp_banned_until <= NOW()) AND otp_timeout_count + 1 >= 2 THEN NOW() + INTERVAL '5 minutes'
				ELSE otp_banned_until
			END
		WHERE tel = ?
		RETURNING otp_timeout_count, otp_banned_until
	`, tel).Scan(ctx, &timeoutCount, &bannedUntil)
	if err != nil {
		return nil, 0, err
	}
	if !bannedUntil.Valid {
		return nil, timeoutCount, nil
	}
	t := bannedUntil.Time
	return &t, timeoutCount, nil
}

func (r user) UpdateProfile(ctx context.Context, p entity.ProfileUpdate) error {
	res, err := r.bun.NewRaw(`
		UPDATE users
		SET tel = ?,
			first_name = ?,
			last_name = ?,
			address_primary = ?,
			address = ?,
			soi = ?,
			road = ?,
			sub_district = ?,
			district = ?,
			province = ?,
			zip_code = ?,
			email = ?,
			facebook = ?,
			updated_at = NOW()
		WHERE user_id = ?
	`, p.Tel, p.FirstName, p.LastName, p.AddressPrimary, p.Address, p.Soi, p.Road,
		p.SubDistrict, p.District, p.Province, p.ZipCode, p.Email, p.Facebook, p.UserID).Exec(ctx)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNoUserUpdated
	}
	return nil
}
