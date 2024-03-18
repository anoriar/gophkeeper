package auth

import (
	"context"

	"github.com/anoriar/gophkeeper/internal/client/user/dto/command"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, command command.RegisterCommand) error
	Login(ctx context.Context, command command.LoginCommand) error
}
