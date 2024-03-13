package entry

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command_response"
	entryFactoryPkg "github.com/anoriar/gophkeeper/internal/client/entry/factory"
	"github.com/anoriar/gophkeeper/internal/client/entry/factory/command/response"
	"github.com/anoriar/gophkeeper/internal/client/entry/factory/ext_repository/request"
	entryRepository "github.com/anoriar/gophkeeper/internal/client/entry/repository/entry"
	"github.com/anoriar/gophkeeper/internal/client/entry/repository/entry_ext"
	"github.com/anoriar/gophkeeper/internal/client/entry/services/encoder"
	sharedErrors "github.com/anoriar/gophkeeper/internal/client/shared/errors"
	"github.com/anoriar/gophkeeper/internal/client/user/repository/secret"
)

type EntryService struct {
	entryFactory       entryFactoryPkg.EntryFactoryInterface
	entryRepository    entryRepository.EntryRepositoryInterface
	secretRepository   secret.SecretRepositoryInterface
	encoder            encoder.DataEncryptorInterface
	responseFactory    *response.EntryResponseFactory
	extEntryRepository entry_ext.EntryExtRepositoryInterface
	syncRequestFactory *request.SyncRequestFactory
	logger             *zap.Logger
}

func NewEntryService(
	entryFactory entryFactoryPkg.EntryFactoryInterface,
	entryRepository entryRepository.EntryRepositoryInterface,
	secretRepository secret.SecretRepositoryInterface,
	encoderInterface encoder.DataEncryptorInterface,
	extEntryRepository entry_ext.EntryExtRepositoryInterface,
	logger *zap.Logger,
) *EntryService {
	return &EntryService{
		entryFactory:       entryFactory,
		entryRepository:    entryRepository,
		secretRepository:   secretRepository,
		encoder:            encoderInterface,
		responseFactory:    response.NewEntryResponseFactory(),
		extEntryRepository: extEntryRepository,
		syncRequestFactory: request.NewSyncRequestFactory(),
		logger:             logger,
	}
}

func (l *EntryService) Add(ctx context.Context, command command.AddEntryCommand) error {
	masterPass, err := l.secretRepository.GetMasterPassword()

	if err != nil {
		if errors.Is(err, secret.ErrMasterPasswordNotFound) {
			return fmt.Errorf("%w: %w", secret.ErrMasterPasswordNotFound, err)
		}
		l.logger.Error("get master password error", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}
	entry, err := l.entryFactory.CreateFromAddCmd(command)
	if err != nil {
		l.logger.Error("create entry error", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}
	encodedData, err := l.encoder.Encrypt(entry.Data, masterPass)
	if err != nil {
		l.logger.Error("encrypt data error", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}
	entry.Data = encodedData
	err = l.entryRepository.Add(ctx, entry)
	if err != nil {
		l.logger.Error("save data error", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}

	return nil
}

func (l *EntryService) Edit(ctx context.Context, command command.EditEntryCommand) error {
	masterPass, err := l.secretRepository.GetMasterPassword()
	if err != nil {
		if errors.Is(err, secret.ErrMasterPasswordNotFound) {
			return fmt.Errorf("%w: %w", secret.ErrMasterPasswordNotFound, err)
		}
		l.logger.Error("get master password error", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}

	entry, err := l.entryFactory.CreateFromEditCmd(command)
	if err != nil {
		l.logger.Error("edit entry error", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}
	encodedData, err := l.encoder.Encrypt(entry.Data, masterPass)
	if err != nil {
		l.logger.Error("encrypt data error", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}
	entry.Data = encodedData
	err = l.entryRepository.Edit(ctx, entry)
	if err != nil {
		if errors.Is(err, sharedErrors.ErrEntryNotFound) {
			return fmt.Errorf("%w: %w", sharedErrors.ErrEntryNotFound, err)
		}
		l.logger.Error("edit data error", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}

	return nil
}

func (l *EntryService) Delete(ctx context.Context, command command.DeleteEntryCommand) error {
	entryEntity, err := l.entryRepository.GetById(ctx, command.Id)
	if err != nil {
		if errors.Is(err, sharedErrors.ErrEntryNotFound) {
			return fmt.Errorf("%w: %w", sharedErrors.ErrEntryNotFound, err)
		}
		l.logger.Error("get entry error", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}
	entryEntity.IsDeleted = true

	err = l.entryRepository.Edit(ctx, entryEntity)
	if err != nil {
		l.logger.Error("delete data error", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}
	return nil
}

func (l *EntryService) Detail(ctx context.Context, command command.DetailEntryCommand) (command_response.DetailEntryCommandResponse, error) {
	masterPass, err := l.secretRepository.GetMasterPassword()
	if err != nil {
		if errors.Is(err, secret.ErrMasterPasswordNotFound) {
			return command_response.DetailEntryCommandResponse{}, fmt.Errorf("%w: %w", secret.ErrMasterPasswordNotFound, err)
		}
		l.logger.Error("get master password error", zap.String("error", err.Error()))
		return command_response.DetailEntryCommandResponse{}, fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}

	entry, err := l.entryRepository.GetById(ctx, command.Id)
	if err != nil {
		l.logger.Error("detail data error", zap.String("error", err.Error()))
		return command_response.DetailEntryCommandResponse{}, fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}
	decryptedData, err := l.encoder.Decrypt(entry.Data, masterPass)
	if err != nil {
		l.logger.Error("decrypt data error", zap.String("error", err.Error()))
		return command_response.DetailEntryCommandResponse{}, fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}
	entry.Data = decryptedData
	return l.responseFactory.CreateDetailResponseFromEntity(entry)
}

func (l *EntryService) List(ctx context.Context) ([]command_response.ListEntryCommandResponse, error) {
	entries, err := l.entryRepository.GetList(ctx)
	if err != nil {
		l.logger.Error("get list data error", zap.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}

	return l.responseFactory.CreateListResponseFromEntity(entries), nil
}

func (l *EntryService) Sync(ctx context.Context, command command.SyncEntryCommand) error {
	token, err := l.secretRepository.GetAuthToken()
	if err != nil {
		if errors.Is(err, secret.ErrTokenNotFound) {
			return fmt.Errorf("%w: %w", secret.ErrTokenNotFound, err)
		}
		l.logger.Error("get token error", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}
	entries, err := l.entryRepository.GetList(ctx)
	if err != nil {
		l.logger.Error("get entries list error", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}
	syncRequestItems := l.syncRequestFactory.CreateSyncRequest(command.EntryType, entries)
	syncResponse, err := l.extEntryRepository.Sync(ctx, token, syncRequestItems)
	if err != nil {
		l.logger.Error("sync entries error", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}
	newEntries, err := l.entryFactory.CreateFromSyncResponse(syncResponse)
	if err != nil {
		l.logger.Error("create from sync response error", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}
	err = l.entryRepository.Rewrite(ctx, newEntries)
	if err != nil {
		l.logger.Error("rewrite entries error", zap.String("error", err.Error()))
		return fmt.Errorf("%w: %w", sharedErrors.ErrInternalError, err)
	}
	return nil
}
