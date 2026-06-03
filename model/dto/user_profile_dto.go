package dto

// UserProfileResponse is returned for GET /users, GET /users/profile, and after PUT /users/profile.
type UserProfileResponse struct {
	UserID            string `json:"user_id"`
	NationalID        string `json:"national_id,omitempty"`
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
	// HasCreditDebt is true when credit is negative (e.g. after a lost top-up dispute).
	HasCreditDebt bool  `json:"has_credit_debt"`
	CreditDebtBaht int64 `json:"credit_debt_baht,omitempty"`
	// WithdrawalBlocked is true while the user must ship as seller or confirm receipt as buyer on a closed auction.
	WithdrawalBlocked     bool   `json:"withdrawal_blocked"`
	WithdrawalBlockReason string `json:"withdrawal_block_reason,omitempty"`
	// PendingSellerShipCount is closed auctions won by someone where seller has not recorded shipment yet.
	PendingSellerShipCount int `json:"pending_seller_ship_count"`
	// AccountRestricted — ถูกจำกัดการใช้งานเต็มรูปแบบ (เติม/ถอน/บิด) จนกว่า restricted_until
	AccountRestricted       bool   `json:"account_restricted"`
	RestrictedUntil         string `json:"restricted_until,omitempty"`
	RestrictedReason        string `json:"restricted_reason,omitempty"`
	// PostingRestricted — ห้ามโพสประมูลใหม่เท่านั้น (จากรายงานที่ admin accept)
	PostingRestricted       bool   `json:"posting_restricted"`
	PostingRestrictedUntil  string `json:"posting_restricted_until,omitempty"`
	PostingRestrictedReason string `json:"posting_restricted_reason,omitempty"`
	ReputationPoints        int64   `json:"reputation_points"` // star points (rating×2, penalties, admin); internal display uses seller_review_avg_rating
	SellerReviewAvgRating   float64 `json:"seller_review_avg_rating"`
	SellerReviewCount       int     `json:"seller_review_count"`
	AppealPending           bool   `json:"appeal_pending"`
	AppealStatus            string `json:"appeal_status,omitempty"`
	UnreadNotificationCount int    `json:"unread_notification_count"`
}

type SubmitRestrictionAppealRequest struct {
	Reason string `json:"reason" validate:"required"`
}

type RestrictionAppealResponse struct {
	AppealID   int64  `json:"appeal_id,omitempty"`
	Status     string `json:"status"`
	Reason     string `json:"reason,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	ResolvedAt string `json:"resolved_at,omitempty"`
	AdminNote  string `json:"admin_note,omitempty"`
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
// When is_first_registration is true, Tel/FirstName/LastName come from tel_verify signup columns.
type OnboardingStatusResponse struct {
	IsFirstRegistration bool   `json:"is_first_registration"`
	Tel               string `json:"tel,omitempty"`
	FirstName         string `json:"first_name,omitempty"`
	LastName          string `json:"last_name,omitempty"`
	Email             string `json:"email,omitempty"`
}

// TelAvailabilityResponse is returned for GET /users/:tel.
type TelAvailabilityResponse struct {
	OK bool `json:"ok"`
}

// OTPTimeoutRequest is the JSON body for POST /otp/timeout.
type OTPTimeoutRequest struct {
	Tel string `json:"tel" validate:"required"`
}
