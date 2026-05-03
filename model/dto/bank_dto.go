package dto

type BankOptionResponse struct {
	BankID   int64  `json:"bank_id"`
	BankCode string `json:"bank_code"`
	NameTH   string `json:"name_th"`
	NameEN   string `json:"name_en"`
}
