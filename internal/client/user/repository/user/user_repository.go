package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	sharedErr "github.com/anoriar/gophkeeper/internal/client/shared/errors"
	"github.com/anoriar/gophkeeper/internal/client/user/dto/repository/request"
	userErr "github.com/anoriar/gophkeeper/internal/client/user/errors"
)

type UserRepository struct {
	client *resty.Client
}

func NewUserRepository(client *resty.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (u UserRepository) Register(request request.RegisterRequest) (token string, err error) {
	body, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	resp, err := u.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("/api/user/register")

	if err != nil {
		return "", err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		token = resp.Header().Get("Authorization")
		return token, nil
	case http.StatusConflict:
		return "", fmt.Errorf("%w: %v", userErr.ErrUserExists, resp.Body())
	default:
		return "", fmt.Errorf("%w: %v", sharedErr.ErrDependencyFailure, resp.Body())
	}
}

func (u UserRepository) Login(request request.LoginRequest) (token string, err error) {
	body, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	resp, err := u.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("/api/user/login")

	if err != nil {
		return "", err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		token = resp.Header().Get("Authorization")
		return token, nil
	case http.StatusUnauthorized:
		return "", fmt.Errorf("%w: %v", userErr.ErrUserUnauthorized, resp.Body())
	default:
		return "", fmt.Errorf("%w: %v", sharedErr.ErrDependencyFailure, resp.Body())
	}
}
