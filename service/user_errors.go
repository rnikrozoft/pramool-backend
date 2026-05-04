package service

import "errors"

var (
	ErrTelRequired       = errors.New("tel is required")
	ErrTelAlreadyUsed    = errors.New("tel is already used")
	ErrUserNotRegistered = errors.New("no users row for this account; complete registration first")

	// Registration / full account (Thai messages for API responses)
	ErrNationalIDAlreadyRegistered = errors.New("หมายเลขบัตรประชาชนนี้ถูกใช้ลงทะเบียนแล้ว")
	ErrEmailAlreadyRegistered      = errors.New("อีเมลนี้ถูกใช้ในระบบแล้ว")
	ErrTelHasFullUserRecord        = errors.New("เบอร์โทรศัพท์นี้ลงทะเบียนแล้ว")
)
