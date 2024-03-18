package sync

import (
	"context"
	"errors"
	"fmt"

	sharedErrors "github.com/anoriar/gophkeeper/internal/server/shared/errors"

	"github.com/anoriar/gophkeeper/internal/server/shared/services/uuid"

	"go.uber.org/zap"

	serverErrors "github.com/anoriar/gophkeeper/internal/server/entry/errors"
	"github.com/anoriar/gophkeeper/internal/server/entry/validator"

	"github.com/anoriar/gophkeeper/internal/server/entry/dto/collection"
	"github.com/anoriar/gophkeeper/internal/server/entry/dto/request/sync"
	syncResponsePkg "github.com/anoriar/gophkeeper/internal/server/entry/dto/response/sync"
	"github.com/anoriar/gophkeeper/internal/server/entry/entity"
	"github.com/anoriar/gophkeeper/internal/server/entry/factory"
	syncResponseFactory "github.com/anoriar/gophkeeper/internal/server/entry/factory/response/sync"
	"github.com/anoriar/gophkeeper/internal/server/entry/repository"
	"github.com/anoriar/gophkeeper/internal/server/shared/app/db"
	context2 "github.com/anoriar/gophkeeper/internal/server/shared/context"
)

type SyncService struct {
	entryRepository      repository.EntryRepositoryInterface
	entryFactory         *factory.EntryFactory
	syncResponseFactory  *syncResponseFactory.SyncResponseFactory
	syncRequestValidator *validator.SyncRequestValidator
	db                   db.DatabaseInterface
	logger               *zap.Logger
}

func NewSyncService(
	entryRepository repository.EntryRepositoryInterface,
	uuidGen uuid.UUIDGeneratorInterface,
	db db.DatabaseInterface,
	logger *zap.Logger,
) *SyncService {
	return &SyncService{
		entryRepository:      entryRepository,
		db:                   db,
		logger:               logger,
		entryFactory:         factory.NewEntryFactory(uuidGen),
		syncResponseFactory:  syncResponseFactory.NewSyncResponseFactory(),
		syncRequestValidator: validator.NewSyncRequestValidator(),
	}
}

func (s SyncService) Sync(ctx context.Context, request sync.SyncRequest) (syncResponsePkg.SyncResponse, error) {
	validationErrors := s.syncRequestValidator.ValidateSyncRequest(request)
	if len(validationErrors) > 0 {
		return syncResponsePkg.SyncResponse{}, fmt.Errorf("%w: %v", serverErrors.ErrSyncRequestNotValid, validationErrors)
	}

	userEntries, err := s.entryRepository.GetEntriesByUserIDAndType(ctx, request.UserID, request.SyncType)
	if err != nil {
		s.logger.Error("get user entries error", zap.String("error", err.Error()))
		return syncResponsePkg.SyncResponse{}, fmt.Errorf("%w: %v", sharedErrors.ErrInternalError, err)
	}

	newEntries := s.getNewItems(request, userEntries)
	updatedEntries := s.getUpdatedItems(request, userEntries)
	deletedIds := s.getDeletedIds(request, userEntries)

	err = s.executeSync(ctx, newEntries, updatedEntries, deletedIds)
	if err != nil {
		if errors.Is(err, sharedErrors.ErrConflict) {
			s.logger.Error("execute sync error", zap.String("error", err.Error()))
			return syncResponsePkg.SyncResponse{}, fmt.Errorf("%w", err)
		}
		s.logger.Error("execute sync error", zap.String("error", err.Error()))
		return syncResponsePkg.SyncResponse{}, fmt.Errorf("%w: %v", sharedErrors.ErrInternalError, err)
	}

	actualEntries, err := s.entryRepository.GetEntriesByUserIDAndType(ctx, request.UserID, request.SyncType)
	if err != nil {
		s.logger.Error("get actual entries error", zap.String("error", err.Error()))
		return syncResponsePkg.SyncResponse{}, fmt.Errorf("%w: %v", sharedErrors.ErrInternalError, err)
	}
	response := s.syncResponseFactory.CreateSyncResponse(actualEntries, request.SyncType)
	return response, nil
}

func (s SyncService) executeSync(ctx context.Context, newEntries []entity.Entry, updatedEntries []entity.Entry, deletedIds []string) error {
	txx, err := s.db.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("create transaction error: %v", err)
	}
	ctx = context.WithValue(ctx, context2.TransactionKey, txx)

	defer txx.Rollback()
	if len(newEntries) > 0 {
		err = s.entryRepository.AddEntries(ctx, newEntries)
		if err != nil {
			return fmt.Errorf("add entries error: %w", err)
		}
	}

	if len(updatedEntries) > 0 {
		err = s.entryRepository.UpdateEntries(ctx, updatedEntries)
		if err != nil {
			return fmt.Errorf("executeSync entries error: %v", err)
		}
	}

	if len(deletedIds) > 0 {
		err = s.entryRepository.DeleteEntries(ctx, deletedIds)
		if err != nil {
			return fmt.Errorf("delete entries error: %v", err)
		}
	}

	err = txx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction error: %v", err)
	}
	s.logger.Info("entries added", zap.Int("count", len(newEntries)))
	s.logger.Info("entries updated", zap.Int("count", len(updatedEntries)))
	s.logger.Info("entries deleted", zap.Int("count", len(deletedIds)))

	return nil
}

func (s SyncService) getNewItems(request sync.SyncRequest, userEntries collection.EntryCollection) []entity.Entry {
	newEntries := make([]entity.Entry, 0, len(request.Items))
	for _, requestItem := range request.Items {
		if requestItem.IsDeleted != true && !userEntries.Contains(requestItem.OriginalId) {
			item := s.entryFactory.CreateNewEntryFromRequestItem(requestItem, request.UserID, request.SyncType)
			newEntries = append(newEntries, item)
		}
	}
	return newEntries
}

func (s SyncService) getUpdatedItems(request sync.SyncRequest, userEntries collection.EntryCollection) []entity.Entry {
	updatedEntries := make([]entity.Entry, 0, len(request.Items))
	for _, requestItem := range request.Items {
		if requestItem.OriginalId == "" || requestItem.IsDeleted == true {
			continue
		}
		userEntry := userEntries.FindByOriginalId(requestItem.OriginalId)
		if userEntry != nil {
			if requestItem.UpdatedAt.After(userEntry.UpdatedAt) {
				item := s.entryFactory.CreateEntryFromRequestItem(userEntry.Id, requestItem, request.UserID, request.SyncType)
				updatedEntries = append(updatedEntries, item)
			}
		}
	}
	return updatedEntries
}

func (s SyncService) getDeletedIds(request sync.SyncRequest, userEntries collection.EntryCollection) []string {
	deletedIds := make([]string, 0, len(request.Items))
	for _, requestItem := range request.Items {
		deletedEntry := userEntries.FindByOriginalId(requestItem.OriginalId)
		if requestItem.IsDeleted == true && deletedEntry != nil {
			deletedIds = append(deletedIds, deletedEntry.Id)
		}
	}
	return deletedIds
}
