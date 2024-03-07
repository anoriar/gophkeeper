package entry

import (
	"io"

	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
	"github.com/anoriar/gophkeeper/internal/client/entry/repository/entry/internal/single_file/reader"
	"github.com/anoriar/gophkeeper/internal/client/entry/repository/entry/internal/single_file/writer"
)

type EntrySingleFileRepository struct {
	fileName string
}

func NewEntrySingleFileRepository(fileName string) *EntrySingleFileRepository {
	return &EntrySingleFileRepository{fileName: fileName}
}

func (e *EntrySingleFileRepository) Add(entry entity.Entry) error {
	fileWriter, err := writer.NewEntryFileWriter(e.fileName)
	if err != nil {
		return err
	}
	defer fileWriter.Close()

	err = fileWriter.WriteEntry(entry)
	if err != nil {
		return err
	}
	return nil
}

func (e *EntrySingleFileRepository) Edit(entry entity.Entry) error {
	//TODO implement me
	panic("implement me")
}

func (e *EntrySingleFileRepository) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

func (e *EntrySingleFileRepository) FindById(id string) (*entity.Entry, error) {
	return e.findOneByCondition(func(entry entity.Entry) bool {
		return entry.Id == id
	})
}

func (e *EntrySingleFileRepository) GetList() ([]entity.Entry, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EntrySingleFileRepository) findOneByCondition(condition func(entry entity.Entry) bool) (*entity.Entry, error) {
	fileReader, err := reader.NewEntryFileReader(e.fileName)
	if err != nil {
		return nil, err
	}

	defer fileReader.Close()

	for {
		entry, err := fileReader.ReadEntry()
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}

		if condition(*entry) {
			return entry, nil
		}
	}
}
