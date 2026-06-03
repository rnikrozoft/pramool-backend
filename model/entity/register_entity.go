package entity

type TelVerify struct {
	Tel             string `db:"tel"`
	SignupFirstName string `db:"signup_first_name"`
	SignupLastName  string `db:"signup_last_name"`
	SignupEmail     string `db:"signup_email"`
}
