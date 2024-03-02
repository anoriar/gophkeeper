package user

import (
	"github.com/google/uuid"

	"github.com/anoriar/gophkeeper/internal/server/user/entity"
)

type UserFactory struct {
}

func NewUserFactory() *UserFactory {
	return &UserFactory{}
}

func (factory *UserFactory) Create(login string, password string, salt string) entity.User {
	return entity.User{
		ID:       uuid.NewString(),
		Login:    login,
		Password: password,
		Salt:     salt,
	}
}
