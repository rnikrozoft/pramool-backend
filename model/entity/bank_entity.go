package entity

type Bank struct {
	BankID       int64  `db:"bank_id"`
	BankCode     string `db:"bank_code"`
	BankNameTH   string `db:"bank_name_th"`
	BankNameEN   string `db:"bank_name_en"`
	DisplayOrder int    `db:"display_order"`
}
