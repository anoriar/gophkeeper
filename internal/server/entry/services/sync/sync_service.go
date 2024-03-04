package sync

import (
	"context"
	"fmt"

	"github.com/anoriar/gophkeeper/internal/server/entry/errors"
	"github.com/anoriar/gophkeeper/internal/server/entry/validator"

	"go.uber.org/zap"

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
	db                   *db.Database
	logger               *zap.Logger
}

func NewSyncService(
	entryRepository repository.EntryRepositoryInterface,
	db *db.Database,
	logger *zap.Logger,
) *SyncService {
	return &SyncService{
		entryRepository:      entryRepository,
		db:                   db,
		logger:               logger,
		entryFactory:         factory.NewEntryFactory(),
		syncResponseFactory:  syncResponseFactory.NewSyncResponseFactory(),
		syncRequestValidator: validator.NewSyncRequestValidator(),
	}
}

func (s SyncService) Sync(ctx context.Context, request sync.SyncRequest) (syncResponsePkg.SyncResponse, error) {
	validationErrors := s.syncRequestValidator.ValidateSyncRequest(request)
	if len(validationErrors) > 0 {
		return syncResponsePkg.SyncResponse{}, fmt.Errorf("%w: errors: %v", errors.ErrSyncRequestNotValid, validationErrors)
	}

	userEntries, err := s.entryRepository.GetEntriesByUserID(ctx, request.UserID)
	if err != nil {
		s.logger.Error("get user entries error", zap.String("error", err.Error()))
		return syncResponsePkg.SyncResponse{}, fmt.Errorf("get user entries error: %v", err)
	}

	newEntries, err := s.getNewItems(request, userEntries)
	if err != nil {
		s.logger.Error("get new entries error", zap.String("error", err.Error()))
		return syncResponsePkg.SyncResponse{}, fmt.Errorf("get new entries error: %v", err)
	}

	updatedEntries, err := s.getUpdatedItems(request, userEntries)
	if err != nil {
		s.logger.Error("get updated entries error", zap.String("error", err.Error()))
		return syncResponsePkg.SyncResponse{}, fmt.Errorf("get updated entries error: %v", err)
	}

	deletedIds := s.getDeletedIds(request, userEntries)

	err = s.executeSync(ctx, newEntries, updatedEntries, deletedIds)
	if err != nil {
		s.logger.Error("execute sync error", zap.String("error", err.Error()))
		return syncResponsePkg.SyncResponse{}, fmt.Errorf("execute sync error: %v", err)
	}

	actualEntries, err := s.entryRepository.GetEntriesByUserID(ctx, request.UserID)
	if err != nil {
		s.logger.Error("get actual entries error", zap.String("error", err.Error()))
		return syncResponsePkg.SyncResponse{}, fmt.Errorf("get actual entries error: %v", err)
	}
	response, err := s.syncResponseFactory.CreateSyncResponse(actualEntries)
	if err != nil {
		s.logger.Error("create response dto error", zap.String("error", err.Error()))
		return syncResponsePkg.SyncResponse{}, fmt.Errorf("create response dto error: %v", err)
	}
	return response, nil
}

func (s SyncService) executeSync(ctx context.Context, newEntries []entity.Entry, updatedEntries []entity.Entry, deletedIds []string) error {
	txx, err := s.db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("create transaction error: %v", err)
	}
	ctx = context.WithValue(ctx, context2.TransactionKey, txx)

	defer txx.Rollback()
	err = s.entryRepository.AddEntries(ctx, newEntries)
	if err != nil {
		return fmt.Errorf("add entries error: %v", err)
	}

	err = s.entryRepository.UpdateEntries(ctx, updatedEntries)
	if err != nil {
		return fmt.Errorf("executeSync entries error: %v", err)
	}

	err = s.entryRepository.DeleteEntries(ctx, deletedIds)
	if err != nil {
		return fmt.Errorf("delete entries error: %v", err)
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

func (s SyncService) getNewItems(request sync.SyncRequest, userEntries collection.EntryCollection) ([]entity.Entry, error) {
	newEntries := make([]entity.Entry, 0, len(request.Items))
	for _, requestItem := range request.Items {
		if !userEntries.Contains(requestItem.Id) && requestItem.IsDeleted != true {
			item, err := s.entryFactory.CreateEntryFromRequestItem(requestItem, request.UserID)
			if err != nil {
				return nil, err
			}
			newEntries = append(newEntries, item)
		}
	}
	return newEntries, nil
}

func (s SyncService) getUpdatedItems(request sync.SyncRequest, userEntries collection.EntryCollection) ([]entity.Entry, error) {
	updatedEntries := make([]entity.Entry, 0, len(request.Items))
	for _, requestItem := range request.Items {
		userEntry := userEntries.FindById(requestItem.Id)
		if userEntry != nil && requestItem.IsDeleted != true {
			if requestItem.UpdatedAt.After(userEntry.UpdatedAt) {
				item, err := s.entryFactory.CreateEntryFromRequestItem(requestItem, request.UserID)
				if err != nil {
					return nil, err
				}
				updatedEntries = append(updatedEntries, item)
			}

		}
	}
	return updatedEntries, nil
}

func (s SyncService) getDeletedIds(request sync.SyncRequest, userEntries collection.EntryCollection) []string {
	deletedIds := make([]string, 0, len(request.Items))
	for _, requestItem := range request.Items {
		if requestItem.IsDeleted == true && userEntries.Contains(requestItem.Id) {
			deletedIds = append(deletedIds, requestItem.Id)
		}
	}
	return deletedIds
}
