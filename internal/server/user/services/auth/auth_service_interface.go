package auth

import (
	"context"

	"github.com/anoriar/gophkeeper/internal/server/user/dto/auth"
	"github.com/anoriar/gophkeeper/internal/server/user/dto/requests/login"
	"github.com/anoriar/gophkeeper/internal/server/user/dto/requests/register"
)

//go:generate mockgen -source=auth_service_interface.go -destination=mock/auth_service.go -package=mock
type AuthServiceInterface interface {
	RegisterUser(ctx context.Context, dto register.RegisterUserRequestDto) (string, error)
	LoginUser(ctx context.Context, dto login.LoginUserRequestDto) (string, error)
	ValidateToken(token string) (auth.UserClaims, error)
}
