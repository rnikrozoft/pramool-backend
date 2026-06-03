package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/rnikrozoft/pramool-core/config"
	"github.com/rnikrozoft/pramool-core/exception"
	"golang.org/x/crypto/bcrypt"
)

const (
	ForgotPasswordStatusOK         = "ok"
	ForgotPasswordStatusNotFound   = "not_found"
	ForgotPasswordStatusNoPassword = "no_password"
)

type PasswordResetService interface {
	CheckEligibility(ctx context.Context, rawTel string) (status string, tel string, err error)
	ResetPassword(ctx context.Context, rawTel, token, pin, plainPassword string) error
}

type passwordReset struct {
	appConfigs  config.AppConfigs
	userService UserService
	otpService  OTPService
}

func NewPasswordResetService(appConfigs config.AppConfigs, userService UserService, otpService OTPService) PasswordResetService {
	return passwordReset{
		appConfigs:  appConfigs,
		userService: userService,
		otpService:  otpService,
	}
}

func (s passwordReset) CheckEligibility(ctx context.Context, rawTel string) (string, string, error) {
	tel, err := s.userService.ResolveLoginIdentifier(ctx, rawTel)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ForgotPasswordStatusNotFound, "", nil
		}
		return "", "", err
	}
	if _, err := s.userService.FindUserIDByTelWithFallback(ctx, tel); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ForgotPasswordStatusNotFound, "", nil
		}
		return "", "", err
	}
	hash, err := s.userService.GetLoginPasswordHash(ctx, tel)
	if err != nil {
		return "", "", err
	}
	if strings.TrimSpace(hash) == "" {
		return ForgotPasswordStatusNoPassword, tel, nil
	}
	return ForgotPasswordStatusOK, tel, nil
}

func (s passwordReset) verifyOTP(ctx context.Context, token, pin string) error {
	token = strings.TrimSpace(token)
	pin = strings.TrimSpace(pin)
	if token == "" || pin == "" {
		return exception.BadRequest(errors.New("กรุณากรอกรหัสยืนยัน"))
	}
	if strings.TrimSpace(s.appConfigs.ThaiBulkSMS.APIKey) == "" {
		if len(pin) < 4 {
			return exception.BadRequest(errors.New("รหัสยืนยันไม่ถูกต้อง"))
		}
		return nil
	}
	if err := s.otpService.Verify(ctx, token, pin); err != nil {
		return exception.BadRequest(errors.New("รหัสยืนยันไม่ถูกต้องหรือหมดอายุ"))
	}
	return nil
}

func (s passwordReset) ResetPassword(ctx context.Context, rawTel, token, pin, plainPassword string) error {
	status, tel, err := s.CheckEligibility(ctx, rawTel)
	if err != nil {
		return err
	}
	switch status {
	case ForgotPasswordStatusNotFound:
		return exception.NotFound(errors.New("ไม่พบบัญชีที่ใช้เบอร์นี้"))
	}
	if err := s.verifyOTP(ctx, token, pin); err != nil {
		return err
	}
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.userService.UpdatePasswordHashByTel(ctx, tel, string(hashBytes))
}
