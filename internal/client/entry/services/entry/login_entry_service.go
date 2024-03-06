package entry

import (
	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
	"github.com/anoriar/gophkeeper/internal/client/entry/factory/entry"
	entryRepository "github.com/anoriar/gophkeeper/internal/client/entry/repository/entry"
)

type LoginEntryService struct {
	loginEntryFactory    entry.EntryFactoryInterface
	loginEntryRepository entryRepository.EntryRepositoryInterface
}

func NewLoginEntryService(loginEntryFactory entry.EntryFactoryInterface, loginEntryRepository entryRepository.EntryRepositoryInterface) *LoginEntryService {
	return &LoginEntryService{loginEntryFactory: loginEntryFactory, loginEntryRepository: loginEntryRepository}
}

func (l *LoginEntryService) Add(command command.AddEntryCommand) error {
	//TODO implement me
	// 1 зашифровать данные в byte
	// 2 фабрика: собрать данные в entry
	// 3 storageProvider.Add
	panic("implement me")
}

func (l *LoginEntryService) Edit(command command.EditEntryCommand) error {
	//TODO implement me
	panic("implement me")
}

func (l *LoginEntryService) Delete(command command.DeleteEntryCommand) error {
	//TODO implement me
	panic("implement me")
}

func (l *LoginEntryService) Detail(command command.DetailEntryCommand) (entity.Entry, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LoginEntryService) List(command command.ListEntryCommand) ([]entity.Entry, error) {
	//TODO implement me
	panic("implement me")
}
