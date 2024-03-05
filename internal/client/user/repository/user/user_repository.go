package user

import (
	"encoding/json"
	"github.com/anoriar/gophkeeper/internal/client/user/dto/repository/request"
	"github.com/go-resty/resty/v2"
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

	token = resp.Header().Get("Authorization")

	return token, nil
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

	token = resp.Header().Get("Authorization")

	return token, nil
}
