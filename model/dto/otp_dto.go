package dto

type RequestOTP struct {
	Tel string `json:"tel" validate:"required,max=10"`
}
