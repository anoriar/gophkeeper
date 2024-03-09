package sync

import (
	"fmt"

	"github.com/anoriar/gophkeeper/internal/server/entry/dto/collection"
	"github.com/anoriar/gophkeeper/internal/server/entry/dto/response/sync"
	"github.com/anoriar/gophkeeper/internal/server/entry/entity"
	errors2 "github.com/anoriar/gophkeeper/internal/server/shared/errors"
)

type SyncResponseFactory struct {
}

func NewSyncResponseFactory() *SyncResponseFactory {
	return &SyncResponseFactory{}
}

func (f *SyncResponseFactory) CreateSyncResponse(entryCollection collection.EntryCollection) (sync.SyncResponse, error) {
	syncResponseItems := make([]sync.SyncResponseItem, 0, len(entryCollection.Entries))
	for _, entry := range entryCollection.Entries {
		responseItem, err := f.CreateSyncResponseItem(entry)
		if err != nil {
			return sync.SyncResponse{}, fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
		}
		syncResponseItems = append(syncResponseItems, responseItem)
	}
	return *sync.NewSyncResponse(syncResponseItems), nil
}

func (f *SyncResponseFactory) CreateSyncResponseItem(entry entity.Entry) (sync.SyncResponseItem, error) {
	return *sync.NewSyncResponseItem(entry.OriginalId, entry.EntryType, entry.UpdatedAt, entry.Data, entry.Meta), nil
}
