package sync

import (
	"context"
	"fmt"
	"github.com/anoriar/gophkeeper/internal/server/shared/app/db"
	context2 "github.com/anoriar/gophkeeper/internal/server/shared/context"
	"github.com/anoriar/gophkeeper/internal/server/storage/dto/collection"
	"github.com/anoriar/gophkeeper/internal/server/storage/dto/request"
	"github.com/anoriar/gophkeeper/internal/server/storage/entity"
	"github.com/anoriar/gophkeeper/internal/server/storage/factory"
	"github.com/anoriar/gophkeeper/internal/server/storage/repository"
	"go.uber.org/zap"
)

type SyncService struct {
	entryRepository repository.EntryRepositoryInterface
	entryFactory    *factory.EntryFactory
	db              *db.Database
	logger          *zap.Logger
}

func NewSyncService(entryRepository repository.EntryRepositoryInterface, db *db.Database, logger *zap.Logger) *SyncService {
	return &SyncService{entryRepository: entryRepository, db: db, logger: logger, entryFactory: factory.NewEntryFactory()}
}

func (s SyncService) Sync(ctx context.Context, request request.SyncRequest) ([]entity.Entry, error) {

	userEntries, err := s.entryRepository.GetEntriesByUserID(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user entries error: %v", err)
	}

	newEntries, err := s.getNewItems(request, userEntries)
	if err != nil {
		return nil, fmt.Errorf("get new entries error: %v", err)
	}

	updatedEntries, err := s.getUpdatedItems(request, userEntries)
	if err != nil {
		return nil, fmt.Errorf("get updated entries error: %v", err)
	}

	deletedIds := s.getDeletedIds(request, userEntries)

	txx, err := s.db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create transaction error: %v", err)
	}
	ctx = context.WithValue(ctx, context2.TransactionKey, txx)

	defer txx.Rollback()
	err = s.entryRepository.AddEntries(ctx, newEntries)
	if err != nil {
		return nil, fmt.Errorf("add entries error: %v", err)
	}
	err = s.entryRepository.UpdateEntries(ctx, updatedEntries)
	if err != nil {
		return nil, fmt.Errorf("update entries error: %v", err)
	}
	err = s.entryRepository.DeleteEntries(ctx, deletedIds)
	if err != nil {
		return nil, fmt.Errorf("delete entries error: %v", err)
	}
	err = txx.Commit()
	if err != nil {
		return nil, fmt.Errorf("commit transaction error: %v", err)
	}

	return nil, nil
}

func (s SyncService) getNewItems(request request.SyncRequest, userEntries collection.EntryCollection) ([]entity.Entry, error) {
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

func (s SyncService) getUpdatedItems(request request.SyncRequest, userEntries collection.EntryCollection) ([]entity.Entry, error) {
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

func (s SyncService) getDeletedIds(request request.SyncRequest, userEntries collection.EntryCollection) []string {
	deletedIds := make([]string, 0, len(request.Items))
	for _, requestItem := range request.Items {
		if requestItem.IsDeleted == true && userEntries.Contains(requestItem.Id) {
			deletedIds = append(deletedIds, requestItem.Id)
		}
	}
	return deletedIds
}
