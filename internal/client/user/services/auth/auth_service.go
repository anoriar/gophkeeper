package auth

import (
	"github.com/anoriar/gophkeeper/internal/client/user/dto/command"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (a *AuthService) Register(command command.RegisterCommand) error {
	//TODO implement me
	panic("implement me")
}

func (a *AuthService) Login(command command.LoginCommand) error {
	//TODO implement me
	panic("implement me")
}
