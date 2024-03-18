package user

import (
	"context"

	"github.com/anoriar/gophkeeper/internal/client/user/dto/repository/request"
)

type UserRepositoryInterface interface {
	Register(ctx context.Context, request request.RegisterRequest) (token string, err error)
	Login(ctx context.Context, request request.LoginRequest) (token string, err error)
}
