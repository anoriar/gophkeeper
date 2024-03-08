package entry

import (
	"context"
	"fmt"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command_response"
	"github.com/anoriar/gophkeeper/internal/client/entry/factory/entry/request"
	"github.com/anoriar/gophkeeper/internal/client/entry/factory/entry/response"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	entryRepository "github.com/anoriar/gophkeeper/internal/client/entry/repository/entry"
	"github.com/anoriar/gophkeeper/internal/client/entry/services/encoder"
	"github.com/anoriar/gophkeeper/internal/client/user/repository/secret"
)

type LoginEntryService struct {
	loginEntryFactory    request.EntryFactoryInterface
	loginEntryRepository entryRepository.EntryRepositoryInterface
	secretRepository     secret.SecretRepositoryInterface
	encoder              encoder.DataEncryptorInterface
	responseFactory      *response.EntryResponseFactory
}

func NewLoginEntryService(
	loginEntryFactory request.EntryFactoryInterface,
	loginEntryRepository entryRepository.EntryRepositoryInterface,
	secretRepository secret.SecretRepositoryInterface,
	encoderInterface encoder.DataEncryptorInterface,
) *LoginEntryService {
	return &LoginEntryService{
		loginEntryFactory:    loginEntryFactory,
		loginEntryRepository: loginEntryRepository,
		secretRepository:     secretRepository,
		encoder:              encoderInterface,
		responseFactory:      response.NewEntryResponseFactory(),
	}
}

func (l *LoginEntryService) Add(ctx context.Context, command command.AddEntryCommand) error {
	masterPass, err := l.secretRepository.GetMasterPassword()

	if err != nil {
		return fmt.Errorf("get masterpass error: %v", err)
	}
	entry, err := l.loginEntryFactory.CreateFromAddCmd(command)
	if err != nil {
		return fmt.Errorf("create entry error: %v", err)
	}
	encodedData, err := l.encoder.Encrypt(entry.Data, masterPass)
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
	masterPass, err := l.secretRepository.GetMasterPassword()

	entry, err := l.loginEntryFactory.CreateFromEditCmd(command)
	if err != nil {
		return fmt.Errorf("edit entry error: %v", err)
	}
	encodedData, err := l.encoder.Encrypt(entry.Data, masterPass)
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

func (l *LoginEntryService) Detail(ctx context.Context, command command.DetailEntryCommand) (command_response.DetailEntryCommandResponse, error) {
	masterPass, err := l.secretRepository.GetMasterPassword()

	entry, err := l.loginEntryRepository.GetById(ctx, command.Id)
	if err != nil {
		return command_response.DetailEntryCommandResponse{}, fmt.Errorf("detail data error: %v", err)
	}
	decryptedData, err := l.encoder.Decrypt(entry.Data, masterPass)
	if err != nil {
		return command_response.DetailEntryCommandResponse{}, fmt.Errorf("decrypt data error: %v", err)
	}
	entry.Data = decryptedData
	return l.responseFactory.CreateDetailResponseFromEntity(entry)
}

func (l *LoginEntryService) List(ctx context.Context, command command.ListEntryCommand) ([]command_response.ListEntryCommandResponse, error) {
	entries, err := l.loginEntryRepository.GetList(ctx)
	if err != nil {
		return nil, fmt.Errorf("get list data error: %v", err)
	}

	return l.responseFactory.CreateListResponseFromEntity(entries), nil
}
