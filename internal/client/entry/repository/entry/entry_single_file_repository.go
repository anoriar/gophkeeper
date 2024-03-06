package entry

import (
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
)

type EntrySingleFileRepository struct {
	fileName string
}

func NewEntrySingleFileRepository(fileName string) *EntrySingleFileRepository {
	return &EntrySingleFileRepository{fileName: fileName}
}

func (e EntrySingleFileRepository) Add(entry entity.Entry) error {
	//TODO implement me
	panic("implement me")
}

func (e EntrySingleFileRepository) Edit(entry entity.Entry) error {
	//TODO implement me
	panic("implement me")
}

func (e EntrySingleFileRepository) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

func (e EntrySingleFileRepository) GetById(id string) (entity.Entry, error) {
	//TODO implement me
	panic("implement me")
}

func (e EntrySingleFileRepository) GetList() ([]entity.Entry, error) {
	//TODO implement me
	panic("implement me")
}
