package entry

import (
	"context"
	"fmt"
	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command_response"
	entryFactoryPkg "github.com/anoriar/gophkeeper/internal/client/entry/factory"
	"github.com/anoriar/gophkeeper/internal/client/entry/factory/command/response"
	"github.com/anoriar/gophkeeper/internal/client/entry/factory/ext_repository/request"
	entryRepository "github.com/anoriar/gophkeeper/internal/client/entry/repository/entry"
	"github.com/anoriar/gophkeeper/internal/client/entry/repository/entry_ext"
	"github.com/anoriar/gophkeeper/internal/client/entry/services/encoder"
	"github.com/anoriar/gophkeeper/internal/client/user/repository/secret"
)

type LoginEntryService struct {
	entryFactory         entryFactoryPkg.EntryFactoryInterface
	loginEntryRepository entryRepository.EntryRepositoryInterface
	secretRepository     secret.SecretRepositoryInterface
	encoder              encoder.DataEncryptorInterface
	responseFactory      *response.EntryResponseFactory
	extEntryRepository   entry_ext.EntryExtRepositoryInterface
	syncRequestFactory   *request.SyncRequestFactory
}

func NewLoginEntryService(
	entryFactory entryFactoryPkg.EntryFactoryInterface,
	loginEntryRepository entryRepository.EntryRepositoryInterface,
	secretRepository secret.SecretRepositoryInterface,
	encoderInterface encoder.DataEncryptorInterface,
	extEntryRepository entry_ext.EntryExtRepositoryInterface,
) *LoginEntryService {
	return &LoginEntryService{
		entryFactory:         entryFactory,
		loginEntryRepository: loginEntryRepository,
		secretRepository:     secretRepository,
		encoder:              encoderInterface,
		responseFactory:      response.NewEntryResponseFactory(),
		extEntryRepository:   extEntryRepository,
		syncRequestFactory:   request.NewSyncRequestFactory(),
	}
}

func (l *LoginEntryService) Add(ctx context.Context, command command.AddEntryCommand) error {
	masterPass, err := l.secretRepository.GetMasterPassword()

	if err != nil {
		return fmt.Errorf("get masterpass error: %w", err)
	}
	entry, err := l.entryFactory.CreateFromAddCmd(command)
	if err != nil {
		return fmt.Errorf("create entry error: %w", err)
	}
	encodedData, err := l.encoder.Encrypt(entry.Data, masterPass)
	if err != nil {
		return fmt.Errorf("encode data error: %w", err)
	}
	entry.Data = encodedData
	err = l.loginEntryRepository.Add(ctx, entry)
	if err != nil {
		return fmt.Errorf("save data error: %w", err)
	}

	return nil
}

func (l *LoginEntryService) Edit(ctx context.Context, command command.EditEntryCommand) error {
	masterPass, err := l.secretRepository.GetMasterPassword()

	entry, err := l.entryFactory.CreateFromEditCmd(command)
	if err != nil {
		return fmt.Errorf("edit entry error: %w", err)
	}
	encodedData, err := l.encoder.Encrypt(entry.Data, masterPass)
	if err != nil {
		return fmt.Errorf("encode data error: %w", err)
	}
	entry.Data = encodedData
	err = l.loginEntryRepository.Edit(ctx, entry)
	if err != nil {
		return fmt.Errorf("edit data error: %w", err)
	}

	return nil
}

func (l *LoginEntryService) Delete(ctx context.Context, command command.DeleteEntryCommand) error {
	entryEntity, err := l.loginEntryRepository.GetById(ctx, command.Id)
	if err != nil {
		return fmt.Errorf("get entry error: %w", err)
	}
	entryEntity.IsDeleted = true

	err = l.loginEntryRepository.Edit(ctx, entryEntity)
	if err != nil {
		return fmt.Errorf("delete data error: %w", err)
	}
	return nil
}

func (l *LoginEntryService) Detail(ctx context.Context, command command.DetailEntryCommand) (command_response.DetailEntryCommandResponse, error) {
	masterPass, err := l.secretRepository.GetMasterPassword()

	entry, err := l.loginEntryRepository.GetById(ctx, command.Id)
	if err != nil {
		return command_response.DetailEntryCommandResponse{}, fmt.Errorf("detail data error: %w", err)
	}
	decryptedData, err := l.encoder.Decrypt(entry.Data, masterPass)
	if err != nil {
		return command_response.DetailEntryCommandResponse{}, fmt.Errorf("decrypt data error: %w", err)
	}
	entry.Data = decryptedData
	return l.responseFactory.CreateDetailResponseFromEntity(entry)
}

func (l *LoginEntryService) List(ctx context.Context) ([]command_response.ListEntryCommandResponse, error) {
	entries, err := l.loginEntryRepository.GetList(ctx)
	if err != nil {
		return nil, fmt.Errorf("get list data error: %w", err)
	}

	return l.responseFactory.CreateListResponseFromEntity(entries), nil
}

func (l *LoginEntryService) Sync(ctx context.Context) error {
	token, err := l.secretRepository.GetAuthToken()
	if err != nil {
		return fmt.Errorf("get token error: %w", err)
	}
	entries, err := l.loginEntryRepository.GetList(ctx)
	if err != nil {
		return fmt.Errorf("get entries list error: %w", err)
	}
	syncRequestItems := l.syncRequestFactory.CreateFromEntries(entries)
	syncResponse, err := l.extEntryRepository.Sync(ctx, token, syncRequestItems)
	if err != nil {
		return fmt.Errorf("sync entries error: %w", err)
	}
	newEntries := l.entryFactory.CreateFromSyncResponse(syncResponse.Items)
	err = l.loginEntryRepository.Rewrite(ctx, newEntries)
	if err != nil {
		return fmt.Errorf("rewrite entries error: %w", err)
	}
	return nil
}
