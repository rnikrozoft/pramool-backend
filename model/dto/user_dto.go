package dto

type UserRegisterRequest struct {
	UserID         string `json:"user_id" validate:"required,max=13"`
	Tel            string `json:"tel" validate:"required,max=10"`
	Email          string `json:"email"`
	Facebook       string `json:"facebook"`
	FirstName      string `json:"first_name" validate:"required,max=100"`
	LastName       string `json:"last_name" validate:"required,max=100"`
	AddressPrimary string `json:"address_primary" validate:"required"`
	Address        string `json:"address"`
	Soi            string `json:"soi"`
	Road           string `json:"road"`
	SubDistrict    string `json:"sub_district" validate:"required"`
	District       string `json:"district" validate:"required"`
	Province       string `json:"province" validate:"required"`
	ZipCode        string `json:"zip_code" validate:"required"`
	Credit         int64  `json:"credit"`
}
