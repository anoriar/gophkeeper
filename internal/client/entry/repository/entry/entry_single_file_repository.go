package entry

import (
	"context"
	"errors"
	"fmt"
	"io"

	sharedErr "github.com/anoriar/gophkeeper/internal/client/shared/errors"

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

func (e *EntrySingleFileRepository) Add(ctx context.Context, entry entity.Entry) error {
	fileWriter, err := writer.NewEntryFileWriter(e.fileName)
	if err != nil {
		return fmt.Errorf("%w: %v", sharedErr.ErrInternalError, err)
	}
	defer fileWriter.Close()

	err = fileWriter.WriteEntry(entry)
	if err != nil {
		return fmt.Errorf("%w: %v", sharedErr.ErrInternalError, err)
	}
	return nil
}

func (e *EntrySingleFileRepository) Edit(ctx context.Context, entry entity.Entry) error {
	return e.rewriteFile(func(fileEntries map[string]*entity.Entry) error {
		if fileEntry, ok := fileEntries[entry.Id]; ok {
			*fileEntry = entry
		} else {
			return fmt.Errorf("%w", sharedErr.ErrEntryNotFound)
		}

		return nil
	})
}

func (e *EntrySingleFileRepository) GetById(ctx context.Context, id string) (entity.Entry, error) {
	entry, err := e.findOneByCondition(func(entry entity.Entry) bool {
		return entry.Id == id
	})
	if err != nil {
		return entity.Entry{}, fmt.Errorf("%w", sharedErr.ErrInternalError)
	}
	if entry == nil {
		return entity.Entry{}, fmt.Errorf("%w", sharedErr.ErrEntryNotFound)
	}
	return *entry, err
}

func (e *EntrySingleFileRepository) GetList(ctx context.Context) ([]entity.Entry, error) {
	fileEntries := make([]entity.Entry, 0)
	fileReader, err := reader.NewEntryFileReader(e.fileName)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", sharedErr.ErrInternalError, err)
	}

	defer fileReader.Close()

	//Считываем все данные с файла
	for {
		entry, err := fileReader.ReadEntry()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("%w: %v", sharedErr.ErrInternalError, err)
		}
		fileEntries = append(fileEntries, *entry)
	}
	return fileEntries, nil
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

func (e *EntrySingleFileRepository) rewriteFile(callback func(fileEntries map[string]*entity.Entry) error) error {
	fileReader, err := reader.NewEntryFileReader(e.fileName)
	if err != nil {
		return nil
	}
	fileEntries := make(map[string]*entity.Entry)

	defer fileReader.Close()

	//Считываем все данные с файла
	for {
		entry, err := fileReader.ReadEntry()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		fileEntries[entry.Id] = entry
	}

	err = callback(fileEntries)
	if err != nil {
		return err
	}

	//Перезаписываем файл заново
	fileWriter, err := writer.NewEntryFileEmptyWriter(e.fileName)
	if err != nil {
		return err
	}
	defer fileWriter.Close()

	for _, entry := range fileEntries {
		err = fileWriter.WriteEntry(*entry)
		if err != nil {
			return err
		}
	}

	return nil
}
