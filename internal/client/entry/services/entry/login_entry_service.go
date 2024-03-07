package entry

import (
	"context"
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

func (l *LoginEntryService) Add(ctx context.Context, command command.AddEntryCommand) error {
	//TODO: разобраться с получением пароля из файла
	//masterPass, err := l.secretRepository.GetMasterPassword()
	//
	//if err != nil {
	//	return fmt.Errorf("get masterpass error: %v", err)
	//}
	//TODO: разобраться с aes. сейчас берет только 16 байт.
	masterPass := "dsah8f3ifnw7f4f4"
	entry, err := l.loginEntryFactory.CreateFromAddCmd(command)
	if err != nil {
		return fmt.Errorf("create entry error: %v", err)
	}
	encodedData, err := l.encoder.Encode(entry.Data, []byte(masterPass))
	if err != nil {
		return fmt.Errorf("encode data error: %v", err)
	}
	entry.Data = encodedData
	err = l.loginEntryRepository.Add(ctx, entry)
	if err != nil {
		return fmt.Errorf("save data error: %v", err)
	}

	return nil
}

func (l *LoginEntryService) Edit(ctx context.Context, command command.EditEntryCommand) error {
	//TODO: masterpass from file + aes
	masterPass := "dsah8f3ifnw7f4f4"
	entry, err := l.loginEntryFactory.CreateFromEditCmd(command)
	if err != nil {
		return fmt.Errorf("edit entry error: %v", err)
	}
	encodedData, err := l.encoder.Encode(entry.Data, []byte(masterPass))
	if err != nil {
		return fmt.Errorf("encode data error: %v", err)
	}
	entry.Data = encodedData
	err = l.loginEntryRepository.Edit(ctx, entry)
	if err != nil {
		return fmt.Errorf("edit data error: %v", err)
	}

	return nil
}

func (l *LoginEntryService) Delete(ctx context.Context, command command.DeleteEntryCommand) error {
	entryEntity, err := l.loginEntryRepository.GetById(ctx, command.Id)
	if err != nil {
		return fmt.Errorf("get entry error: %v", err)
	}
	entryEntity.IsDeleted = true

	err = l.loginEntryRepository.Edit(ctx, entryEntity)
	if err != nil {
		return fmt.Errorf("delete data error: %v", err)
	}
	return nil
}

func (l *LoginEntryService) Detail(ctx context.Context, command command.DetailEntryCommand) (entity.Entry, error) {
	//TODO implement me
	// 1 Получить по ид запись
	// 2 получить мастер пароль
	// 3 расшифровать по aes
	// 4 преобразовать в ответ (как и реквест)
	panic("implement me")
}

func (l *LoginEntryService) List(ctx context.Context, command command.ListEntryCommand) ([]entity.Entry, error) {
	//TODO implement me
	panic("implement me")
}
