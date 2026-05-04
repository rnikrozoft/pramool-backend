package entity

import "time"

type User struct {
	UserID            string    `db:"user_id"`
	Tel               string    `db:"tel"`
	Email             string    `db:"email"`
	Facebook          string    `db:"facebook"`
	BankID            int64     `db:"bank_id"`
	BankAccountName   string    `db:"bank_account_name"`
	BankAccountNumber string    `db:"bank_account_number"`
	FirstName         string    `db:"first_name"`
	LastName          string    `db:"last_name"`
	AddressPrimary    string    `db:"address_primary"`
	Address           string    `db:"address"`
	Soi               string    `db:"soi"`
	Road              string    `db:"road"`
	SubDistrict       string    `db:"sub_district"`
	District          string    `db:"district"`
	Province          string    `db:"province"`
	ZipCode           string    `db:"zip_code"`
	Credit            int64     `db:"credit"`
	PasswordHash      string    `db:"password_hash"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}
