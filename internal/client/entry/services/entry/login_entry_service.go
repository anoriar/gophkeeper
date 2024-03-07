package entry

import (
	"fmt"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
	"github.com/anoriar/gophkeeper/internal/client/entry/factory/entry"
	entryRepository "github.com/anoriar/gophkeeper/internal/client/entry/repository/entry"
	"github.com/anoriar/gophkeeper/internal/client/entry/services/encoder"
	"github.com/anoriar/gophkeeper/internal/client/user/repository/secret"
)

type LoginEntryService struct {
	loginEntryFactory    entry.EntryFactoryInterface
	loginEntryRepository entryRepository.EntryRepositoryInterface
	secretRepository     secret.SecretRepositoryInterface
	encoder              encoder.DataEncoderInterface
}

func NewLoginEntryService(
	loginEntryFactory entry.EntryFactoryInterface,
	loginEntryRepository entryRepository.EntryRepositoryInterface,
	secretRepository secret.SecretRepositoryInterface,
	encoderInterface encoder.DataEncoderInterface,
) *LoginEntryService {
	return &LoginEntryService{
		loginEntryFactory:    loginEntryFactory,
		loginEntryRepository: loginEntryRepository,
		secretRepository:     secretRepository,
		encoder:              encoderInterface,
	}
}

func (l *LoginEntryService) Add(command command.AddEntryCommand) error {

	masterPass, err := l.secretRepository.GetMasterPassword()
	if err != nil {
		return fmt.Errorf("get masterpass error: %v", err)
	}
	entry, err := l.loginEntryFactory.CreateFromAddCmd(command)
	if err != nil {
		return fmt.Errorf("create entry error: %v", err)
	}
	encodedData, err := l.encoder.Encode(entry.Data, []byte(masterPass))
	if err != nil {
		return fmt.Errorf("encode data error: %v", err)
	}
	entry.Data = encodedData
	err = l.loginEntryRepository.Add(entry)
	if err != nil {
		return fmt.Errorf("save data error: %v", err)
	}

	return nil
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
