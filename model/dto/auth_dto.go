package dto

// LoginTelRequest is the JSON body for POST /login/tel.
// Set either Login (recommended: phone or email) or Tel (phone only, backward compatible).
// If the account has a password set, Password is required.
type LoginTelRequest struct {
	Login    string `json:"login" validate:"required_without=Tel,omitempty,max=120"`
	Tel      string `json:"tel" validate:"required_without=Login,omitempty,max=10"`
	Password string `json:"password" validate:"omitempty,max=72"`
}

// SignupRequest is the JSON body for POST /auth/signup.
type SignupRequest struct {
	FirstName       string `json:"first_name" validate:"required,max=20"`
	LastName        string `json:"last_name" validate:"required,max=20"`
	Tel             string `json:"tel" validate:"required,max=10"`
	Email           string `json:"email" validate:"omitempty,max=120,email"`
	Password        string `json:"password" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}
