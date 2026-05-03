package dto

// UserProfileResponse is returned for GET /users, GET /users/profile, and after PUT /users/profile.
type UserProfileResponse struct {
	UserID            string `json:"user_id"`
	Tel               string `json:"tel"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	AddressPrimary    string `json:"address_primary"`
	Address           string `json:"address"`
	Soi               string `json:"soi"`
	Road              string `json:"road"`
	SubDistrict       string `json:"sub_district"`
	District          string `json:"district"`
	Province          string `json:"province"`
	ZipCode           string `json:"zip_code"`
	Email             string `json:"email"`
	Facebook          string `json:"facebook"`
	BankID            int64  `json:"bank_id"`
	BankAccountName   string `json:"bank_account_name"`
	BankAccountNumber string `json:"bank_account_number"`
	Credit            int64  `json:"credit"`
	// WithdrawalBlocked is true while the user has pending seller ship or buyer confirm (escrow) on a closed auction.
	WithdrawalBlocked     bool   `json:"withdrawal_blocked"`
	WithdrawalBlockReason string `json:"withdrawal_block_reason,omitempty"`
}

// UpdateProfileRequest is the JSON body for PUT /users/profile.
type UpdateProfileRequest struct {
	Tel               string `json:"tel" validate:"required"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	AddressPrimary    string `json:"address_primary"`
	Address           string `json:"address"`
	Soi               string `json:"soi"`
	Road              string `json:"road"`
	SubDistrict       string `json:"sub_district"`
	District          string `json:"district"`
	Province          string `json:"province"`
	ZipCode           string `json:"zip_code"`
	Email             string `json:"email"`
	Facebook          string `json:"facebook"`
	BankID            int64  `json:"bank_id"`
	BankAccountName   string `json:"bank_account_name"`
	BankAccountNumber string `json:"bank_account_number"`
}

// OnboardingStatusResponse is returned for GET /users/onboarding-status.
type OnboardingStatusResponse struct {
	IsFirstRegistration bool `json:"is_first_registration"`
}

// TelAvailabilityResponse is returned for GET /users/:tel.
type TelAvailabilityResponse struct {
	OK bool `json:"ok"`
}

// OTPTimeoutRequest is the JSON body for POST /otp/timeout.
type OTPTimeoutRequest struct {
	Tel string `json:"tel" validate:"required"`
}
