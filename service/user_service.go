package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/rnikrozoft/pramool-core/repository"
)

type UserService interface {
	IsTelAlreadyUsed(ctx context.Context, tel string) (bool, error)
	GetMyInfo(ctx context.Context, userID string) (*entity.User, error)
	CountUserFulfillmentBlocks(ctx context.Context, userID string) (pendingSellerShip int, pendingBuyerConfirm int, err error)
	FindUserIDByTelWithFallback(ctx context.Context, tel string) (string, error)
	IsFirstRegistration(ctx context.Context, subject string) (bool, error)
	UpdateMyProfile(ctx context.Context, p entity.ProfileUpdate) error
	SaveUser(ctx context.Context, u *entity.User) error
	GetActiveOTPBanUntil(ctx context.Context, tel string) (*time.Time, error)
	RecordOTPTimeout(ctx context.Context, tel string) (*time.Time, int, error)
}

type user struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return user{
		userRepository: userRepository,
	}
}

func (s user) CountUserFulfillmentBlocks(ctx context.Context, userID string) (int, int, error) {
	return s.userRepository.CountUserFulfillmentBlocks(ctx, strings.TrimSpace(userID))
}

func (s user) IsTelAlreadyUsed(ctx context.Context, tel string) (bool, error) {
	return s.userRepository.IsTelAlreadyUsed(ctx, tel)
}

func (s user) loadUserOrTelVerify(ctx context.Context, userID string) (*entity.User, error) {
	u, err := s.userRepository.FindByID(ctx, userID)
	if err == nil {
		return u, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	tv, err := s.userRepository.FindTelVerifyByTel(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &entity.User{UserID: userID, Tel: tv.Tel}, nil
}

func (s user) GetMyInfo(ctx context.Context, userID string) (*entity.User, error) {
	u, err := s.loadUserOrTelVerify(ctx, userID)
	if err == nil {
		return u, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	realUserID, findErr := s.userRepository.FindUserIDByTel(ctx, strings.TrimSpace(userID))
	if findErr == nil && strings.TrimSpace(realUserID) != "" {
		return s.loadUserOrTelVerify(ctx, strings.TrimSpace(realUserID))
	}
	if findErr != nil && !errors.Is(findErr, sql.ErrNoRows) {
		return nil, findErr
	}

	return s.loadUserOrTelVerify(ctx, userID)
}

func (s user) FindUserIDByTelWithFallback(ctx context.Context, tel string) (string, error) {
	userID, err := s.userRepository.FindUserIDByTel(ctx, tel)
	if err == nil {
		return userID, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	return tel, nil
}

func (s user) IsFirstRegistration(ctx context.Context, subject string) (bool, error) {
	hasUserRecord, err := s.userRepository.HasUserRecordForSubject(ctx, strings.TrimSpace(subject))
	if err != nil {
		return false, err
	}
	return !hasUserRecord, nil
}

func (s user) UpdateMyProfile(ctx context.Context, p entity.ProfileUpdate) error {
	tel := strings.TrimSpace(p.Tel)
	if tel == "" {
		return ErrTelRequired
	}
	usedByOther, err := s.userRepository.IsTelUsedByOtherUser(ctx, p.UserID, tel)
	if err != nil {
		return err
	}
	if usedByOther {
		return ErrTelAlreadyUsed
	}
	err = s.userRepository.UpdateProfile(ctx, entity.ProfileUpdate{
		UserID:            p.UserID,
		Tel:               tel,
		FirstName:         strings.TrimSpace(p.FirstName),
		LastName:          strings.TrimSpace(p.LastName),
		AddressPrimary:    strings.TrimSpace(p.AddressPrimary),
		Address:           strings.TrimSpace(p.Address),
		Soi:               strings.TrimSpace(p.Soi),
		Road:              strings.TrimSpace(p.Road),
		SubDistrict:       strings.TrimSpace(p.SubDistrict),
		District:          strings.TrimSpace(p.District),
		Province:          strings.TrimSpace(p.Province),
		ZipCode:           strings.TrimSpace(p.ZipCode),
		Email:             strings.TrimSpace(p.Email),
		Facebook:          strings.TrimSpace(p.Facebook),
		BankID:            p.BankID,
		BankAccountName:   strings.TrimSpace(p.BankAccountName),
		BankAccountNumber: strings.TrimSpace(p.BankAccountNumber),
	})
	if err != nil {
		if errors.Is(err, repository.ErrNoUserUpdated) {
			return ErrUserNotRegistered
		}
		return err
	}
	return nil
}

func (s user) SaveUser(ctx context.Context, u *entity.User) error {
	if u == nil {
		return errors.New("nil user")
	}
	u.UserID = strings.TrimSpace(u.UserID)
	u.Tel = strings.TrimSpace(u.Tel)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.AddressPrimary = strings.TrimSpace(u.AddressPrimary)
	u.Address = strings.TrimSpace(u.Address)
	u.Soi = strings.TrimSpace(u.Soi)
	u.Road = strings.TrimSpace(u.Road)
	u.SubDistrict = strings.TrimSpace(u.SubDistrict)
	u.District = strings.TrimSpace(u.District)
	u.Province = strings.TrimSpace(u.Province)
	u.ZipCode = strings.TrimSpace(u.ZipCode)
	u.Email = strings.TrimSpace(u.Email)
	u.Facebook = strings.TrimSpace(u.Facebook)
	return s.userRepository.Upsert(ctx, u)
}

func (s user) GetActiveOTPBanUntil(ctx context.Context, tel string) (*time.Time, error) {
	tel = strings.TrimSpace(tel)
	if err := s.userRepository.ResetExpiredOTPBanIfNeeded(ctx, tel); err != nil {
		return nil, err
	}
	return s.userRepository.SelectOTPBanUntil(ctx, tel)
}

func (s user) RecordOTPTimeout(ctx context.Context, tel string) (*time.Time, int, error) {
	return s.userRepository.RecordOTPTimeout(ctx, strings.TrimSpace(tel))
}
