package model

type VerifyRequest struct {
	Token string `json:"token"`
	PIN   string `json:"pin"`
}

type OTPResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
	RefNo  string `json:"refNo"`
}
