package entity

import "time"

type User struct {
	UserID         string    `db:"user_id"`
	Tel            string    `db:"tel"`
	FirstName      string    `db:"first_name"`
	LastName       string    `db:"last_name"`
	AddressPrimary string    `db:"address_primary"`
	Address        string    `db:"address"`
	Soi            string    `db:"soi"`
	Road           string    `db:"road"`
	SubDistrict    string    `db:"sub_district"`
	District       string    `db:"district"`
	Province       string    `db:"province"`
	ZipCode        string    `db:"zip_code"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
