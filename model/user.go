package model

type User struct {
	// UserId    string `json:"user_id" validate:"required,max=13"`
	// FirstName string `json:"first_name" validate:"required,max=20"`
	// LastName  string `json:"last_name" validate:"required,max=20"`
	Email    string `json:"email" validate:"required,max=255"`
	Password string `json:"password" validate:"required,max=255"`
}
