package command

import (
	"fmt"
	validation "github.com/anoriar/gophkeeper/internal/client/shared/dto"
)

type RegisterCommand struct {
	UserName       string
	Password       string
	MasterPassword string
}

func (command *RegisterCommand) Validate() validation.ValidationErrors {
	var validationErrors validation.ValidationErrors
	if command.UserName == "" {
		validationErrors = append(validationErrors, fmt.Errorf("username required"))
	}
	if command.Password == "" {
		validationErrors = append(validationErrors, fmt.Errorf("password required"))
	}

	if command.MasterPassword == "" {
		validationErrors = append(validationErrors, fmt.Errorf("master password required"))
	}

	return validationErrors
}
