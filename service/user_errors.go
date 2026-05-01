package service

import "errors"

var (
	ErrTelRequired       = errors.New("tel is required")
	ErrTelAlreadyUsed    = errors.New("tel is already used")
	ErrUserNotRegistered = errors.New("no users row for this account; complete registration first")
)
