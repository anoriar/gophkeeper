package internal

import (
	"errors"

	validation "github.com/anoriar/gophkeeper/internal/server/shared/dto"
	"github.com/anoriar/gophkeeper/internal/server/user/dto/requests/login"
)

type LoginValidator struct {
}

func NewLoginValidator() *LoginValidator {
	return &LoginValidator{}
}

func (validator *LoginValidator) Validate(requestDto login.LoginUserRequestDto) validation.ValidationErrors {
	var validationErrors validation.ValidationErrors

	if requestDto.Login == "" {
		validationErrors = append(validationErrors, errors.New("login required"))
	}

	if requestDto.Password == "" {
		validationErrors = append(validationErrors, errors.New("password required"))
	}

	return validationErrors
}
