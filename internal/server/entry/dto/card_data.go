package dto

type CardData struct {
	Number     string `json:"number"`
	ExpireDate string `json:"expireDate"`
	Holder     string `json:"holder"`
	CVV        string `json:"cvv"`
}
