package dto

import (
	"fmt"
	validation "github.com/anoriar/gophkeeper/internal/client/shared/dto"
)

type CardData struct {
	Number     string `json:"number"`
	ExpireDate string `json:"expireDate"`
	Holder     string `json:"holder"`
	CVV        string `json:"cvv"`
}

func (data *CardData) Validate() validation.ValidationErrors {
	var validationErrors validation.ValidationErrors
	if data.Number == "" {
		validationErrors = append(validationErrors, fmt.Errorf("number required"))
	}
	if data.Holder == "" {
		validationErrors = append(validationErrors, fmt.Errorf("holder required"))
	}
	if data.ExpireDate == "" {
		validationErrors = append(validationErrors, fmt.Errorf("expire date required"))
	}
	if data.CVV == "" {
		validationErrors = append(validationErrors, fmt.Errorf("cvv required"))
	}

	return validationErrors
}
