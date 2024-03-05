package auth

import "github.com/anoriar/gophkeeper/internal/client/user/dto/command"

type AuthServiceInterface interface {
	Register(command command.RegisterCommand) error
	Login(command command.LoginCommand) error
}
