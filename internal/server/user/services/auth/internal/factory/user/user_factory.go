package user

import (
	"github.com/anoriar/gophkeeper/internal/server/shared/services/uuid"
	"github.com/anoriar/gophkeeper/internal/server/user/entity"
)

type UserFactory struct {
	uuidGen uuid.UUIDGeneratorInterface
}

func NewUserFactory(uuidGen uuid.UUIDGeneratorInterface) *UserFactory {
	return &UserFactory{uuidGen: uuidGen}
}

func (factory *UserFactory) Create(login string, password string, salt string) entity.User {
	return entity.User{
		ID:       factory.uuidGen.NewString(),
		Login:    login,
		Password: password,
		Salt:     salt,
	}
}
