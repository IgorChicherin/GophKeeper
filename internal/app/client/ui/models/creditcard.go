package models

type CreditCard struct {
	Number uint   `json:"number"`
	Month  uint   `json:"month"`
	Year   uint   `json:"year"`
	CVV    uint   `json:"CVV"`
	Holder string `json:"holder"`
}
