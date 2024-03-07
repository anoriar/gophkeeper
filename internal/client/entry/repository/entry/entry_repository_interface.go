package entry

import (
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
)

type EntryRepositoryInterface interface {
	Add(entry entity.Entry) error
	Edit(entry entity.Entry) error
	Delete(id string) error
	FindById(id string) (*entity.Entry, error)
	GetList() ([]entity.Entry, error)
}
