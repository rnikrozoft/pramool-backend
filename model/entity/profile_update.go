package entity

// ProfileUpdate is used for PUT /users/profile.
type ProfileUpdate struct {
	UserID            string
	Tel               string
	FirstName         string
	LastName          string
	AddressPrimary    string
	Address           string
	Soi               string
	Road              string
	SubDistrict       string
	District          string
	Province          string
	ZipCode           string
	Email             string
	Facebook          string
	BankID            int64
	BankAccountName   string
	BankAccountNumber string
}
