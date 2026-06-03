package dto

// LoginTelRequest is the JSON body for POST /login/tel.
// Set either Login (recommended: phone or email) or Tel (phone only, backward compatible).
type LoginTelRequest struct {
	Login    string `json:"login" validate:"required_without=Tel,omitempty,max=120"`
	Tel      string `json:"tel" validate:"required_without=Login,omitempty,max=10"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	Remember bool   `json:"remember"`
}

// ForgotPasswordCheckRequest is the JSON body for POST /auth/forgot-password/check.
type ForgotPasswordCheckRequest struct {
	Tel string `json:"tel" validate:"required,max=10"`
}

// ForgotPasswordResetRequest is the JSON body for POST /auth/forgot-password/reset.
type ForgotPasswordResetRequest struct {
	Tel             string `json:"tel" validate:"required,max=10"`
	Token           string `json:"token" validate:"required"`
	PIN             string `json:"pin" validate:"required"`
	Password        string `json:"password" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

// SignupRequest is the JSON body for POST /auth/signup.
type SignupRequest struct {
	ConsentPayload
	FirstName       string `json:"first_name" validate:"required,max=20"`
	LastName        string `json:"last_name" validate:"required,max=20"`
	Tel             string `json:"tel" validate:"required,max=10"`
	Email           string `json:"email" validate:"omitempty,max=120,email"`
	Password        string `json:"password" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}
