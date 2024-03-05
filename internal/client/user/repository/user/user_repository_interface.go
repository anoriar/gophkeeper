package user

import "github.com/anoriar/gophkeeper/internal/client/user/dto/repository/request"

type UserRepositoryInterface interface {
	Register(request request.RegisterRequest) (token string, err error)
	Login(request request.LoginRequest) (token string, err error)
}
