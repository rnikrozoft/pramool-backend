package entity

type TelVerify struct {
	Tel    string `db:"tel"`
	Verify bool   `db:"verify"`
}
