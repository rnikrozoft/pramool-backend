package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/rnikrozoft/pramool-core/model/dto"
	"github.com/rnikrozoft/pramool-core/model/entity"
	"github.com/rnikrozoft/pramool-core/repository"
)

type UserService interface {
	IsTelAlreadyUsed(ctx context.Context, tel string) (bool, error)
	ExistsRegisteredNationalID(ctx context.Context, nationalID string) (bool, error)
	IsEmailTakenByOtherTel(ctx context.Context, email, requestTel string) (bool, error)
	GetMyInfo(ctx context.Context, userID string) (*entity.User, error)
	CountUserFulfillmentBlocks(ctx context.Context, userID string) (pendingSellerShip int, err error)
	FindUserIDByTelWithFallback(ctx context.Context, tel string) (string, error)
	IsFirstRegistration(ctx context.Context, subject string) (bool, error)
	GetOnboardingStatus(ctx context.Context, subject string) (dto.OnboardingStatusResponse, error)
	UpdateMyProfile(ctx context.Context, p entity.ProfileUpdate) error
	SaveUser(ctx context.Context, u *entity.User) error
	GetActiveOTPBanUntil(ctx context.Context, tel string) (*time.Time, error)
	RecordOTPTimeout(ctx context.Context, tel string) (*time.Time, int, error)
	GetLoginPasswordHash(ctx context.Context, tel string) (string, error)
	UpdatePasswordHashByTel(ctx context.Context, tel, passwordHash string) error
	ResolveLoginIdentifier(ctx context.Context, raw string) (tel string, err error)
	IsUserSuspended(ctx context.Context, subject string) (bool, error)
	SubmitRestrictionAppeal(ctx context.Context, userID, reason string) (int64, error)
	GetRestrictionAppeal(ctx context.Context, userID string) (*repository.RestrictionAppealRow, error)
	HasPendingRestrictionAppeal(ctx context.Context, userID string) (bool, error)
	ProfileAppealMeta(ctx context.Context, userID string) (pending bool, lastStatus string, err error)
	CountUnreadNotifications(ctx context.Context, userID string) (int, error)
	ListNotifications(ctx context.Context, userID string, limit, offset int) ([]repository.UserNotificationRow, int64, error)
	MarkNotificationRead(ctx context.Context, userID string, notificationID int64) (*repository.UserNotificationRow, error)
}

type user struct {
	userRepository         repository.UserRepository
	notificationRepository repository.UserNotificationRepository
}

func NewUserService(userRepository repository.UserRepository, notificationRepository repository.UserNotificationRepository) UserService {
	return user{
		userRepository:         userRepository,
		notificationRepository: notificationRepository,
	}
}

func (s user) CountUserFulfillmentBlocks(ctx context.Context, userID string) (int, error) {
	userID = strings.TrimSpace(userID)
	if !repository.SubjectIsUserUUID(userID) {
		return 0, nil
	}
	return s.userRepository.CountUserFulfillmentBlocks(ctx, userID)
}

func (s user) IsTelAlreadyUsed(ctx context.Context, tel string) (bool, error) {
	return s.userRepository.IsTelAlreadyUsed(ctx, tel)
}

func (s user) ExistsRegisteredNationalID(ctx context.Context, nationalID string) (bool, error) {
	return s.userRepository.ExistsRegisteredNationalID(ctx, nationalID)
}

func (s user) IsEmailTakenByOtherTel(ctx context.Context, email, requestTel string) (bool, error) {
	return s.userRepository.IsEmailTakenByOtherTel(ctx, email, requestTel)
}

func (s user) loadUserOrTelVerify(ctx context.Context, userID string) (*entity.User, error) {
	userID = strings.TrimSpace(userID)
	if repository.SubjectIsUserUUID(userID) {
		u, err := s.userRepository.FindByID(ctx, userID)
		if err == nil {
			return u, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	tv, err := s.userRepository.FindTelVerifyByTel(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &entity.User{
		UserID:    userID,
		Tel:       tv.Tel,
		FirstName: tv.SignupFirstName,
		LastName:  tv.SignupLastName,
	}, nil
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
	// Onboarding login uses tel as JWT subject only after OTP created tel_verify.
	if _, err := s.userRepository.FindTelVerifyByTel(ctx, tel); err != nil {
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

func (s user) GetOnboardingStatus(ctx context.Context, subject string) (dto.OnboardingStatusResponse, error) {
	subject = strings.TrimSpace(subject)
	isFirst, err := s.IsFirstRegistration(ctx, subject)
	if err != nil {
		return dto.OnboardingStatusResponse{}, err
	}
	resp := dto.OnboardingStatusResponse{IsFirstRegistration: isFirst}
	if !isFirst || repository.SubjectIsUserUUID(subject) {
		return resp, nil
	}
	tv, err := s.userRepository.FindTelVerifyByTel(ctx, subject)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			resp.Tel = subject
			return resp, nil
		}
		return dto.OnboardingStatusResponse{}, err
	}
	resp.Tel = strings.TrimSpace(tv.Tel)
	if resp.Tel == "" {
		resp.Tel = subject
	}
	resp.FirstName = strings.TrimSpace(tv.SignupFirstName)
	resp.LastName = strings.TrimSpace(tv.SignupLastName)
	resp.Email = strings.TrimSpace(tv.SignupEmail)
	return resp, nil
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

func (s user) GetLoginPasswordHash(ctx context.Context, tel string) (string, error) {
	return s.userRepository.GetLoginPasswordHash(ctx, strings.TrimSpace(tel))
}

func (s user) UpdatePasswordHashByTel(ctx context.Context, tel, passwordHash string) error {
	return s.userRepository.UpdatePasswordHashByTel(ctx, strings.TrimSpace(tel), strings.TrimSpace(passwordHash))
}

func normalizeThaiLocalTel(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	d := b.String()
	if len(d) >= 10 {
		return d[len(d)-10:]
	}
	if len(d) == 9 {
		return "0" + d
	}
	if len(d) > 0 && len(d) < 9 {
		return d
	}
	return d
}

func (s user) ResolveLoginIdentifier(ctx context.Context, raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", sql.ErrNoRows
	}
	if strings.Contains(raw, "@") {
		return s.userRepository.FindTelByEmail(ctx, raw)
	}
	tel := normalizeThaiLocalTel(raw)
	if tel == "" || len(tel) < 9 {
		return "", sql.ErrNoRows
	}
	if len(tel) > 10 {
		tel = tel[len(tel)-10:]
	}
	return tel, nil
}

func (s user) IsUserSuspended(ctx context.Context, subject string) (bool, error) {
	return s.userRepository.IsUserSuspended(ctx, subject)
}

func (s user) SubmitRestrictionAppeal(ctx context.Context, userID, reason string) (int64, error) {
	reason = strings.TrimSpace(reason)
	if len([]rune(reason)) < 10 {
		return 0, errors.New("reason must be at least 10 characters")
	}
	if len([]rune(reason)) > 2000 {
		return 0, errors.New("reason too long")
	}
	u, err := s.GetMyInfo(ctx, userID)
	if err != nil {
		return 0, err
	}
	return s.userRepository.InsertRestrictionAppeal(ctx, u.UserID, reason)
}

func (s user) GetRestrictionAppeal(ctx context.Context, userID string) (*repository.RestrictionAppealRow, error) {
	u, err := s.GetMyInfo(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.userRepository.GetRestrictionAppealForUser(ctx, u.UserID)
}

func (s user) HasPendingRestrictionAppeal(ctx context.Context, userID string) (bool, error) {
	userID = strings.TrimSpace(userID)
	if !repository.SubjectIsUserUUID(userID) {
		return false, nil
	}
	u, err := s.GetMyInfo(ctx, userID)
	if err != nil {
		return false, err
	}
	return s.userRepository.HasPendingRestrictionAppeal(ctx, u.UserID)
}

func (s user) ProfileAppealMeta(ctx context.Context, userID string) (bool, string, error) {
	userID = strings.TrimSpace(userID)
	if !repository.SubjectIsUserUUID(userID) {
		return false, "", nil
	}
	u, err := s.GetMyInfo(ctx, userID)
	if err != nil {
		return false, "", err
	}
	pending, err := s.userRepository.HasPendingRestrictionAppeal(ctx, u.UserID)
	if err != nil {
		return false, "", err
	}
	if pending {
		return true, "", nil
	}
	row, err := s.userRepository.GetRestrictionAppealForUser(ctx, u.UserID)
	if err != nil {
		return false, "", err
	}
	if row == nil || row.Status == "pending" {
		return false, "", nil
	}
	return false, row.Status, nil
}

func (s user) CountUnreadNotifications(ctx context.Context, userID string) (int, error) {
	userID = strings.TrimSpace(userID)
	if !repository.SubjectIsUserUUID(userID) {
		return 0, nil
	}
	u, err := s.GetMyInfo(ctx, userID)
	if err != nil {
		return 0, err
	}
	return s.notificationRepository.CountUnread(ctx, u.UserID)
}

func (s user) ListNotifications(ctx context.Context, userID string, limit, offset int) ([]repository.UserNotificationRow, int64, error) {
	u, err := s.GetMyInfo(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	return s.notificationRepository.List(ctx, u.UserID, limit, offset)
}

func (s user) MarkNotificationRead(ctx context.Context, userID string, notificationID int64) (*repository.UserNotificationRow, error) {
	u, err := s.GetMyInfo(ctx, userID)
	if err != nil {
		return nil, err
	}
	row, err := s.notificationRepository.MarkRead(ctx, u.UserID, notificationID)
	if err != nil {
		return nil, err
	}
	return row, nil
}
