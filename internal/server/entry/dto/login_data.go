package dto

import (
	"fmt"

	validation "github.com/anoriar/gophkeeper/internal/server/shared/dto"
)

type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (data *LoginData) Validate() validation.ValidationErrors {
	var validationErrors validation.ValidationErrors
	if data.Login == "" {
		validationErrors = append(validationErrors, fmt.Errorf("login required"))
	}
	if data.Password == "" {
		validationErrors = append(validationErrors, fmt.Errorf("password required"))
	}

	return validationErrors
}
