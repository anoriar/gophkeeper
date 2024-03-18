package sync

import (
	"encoding/base64"

	"github.com/anoriar/gophkeeper/internal/server/entry/enum"

	"github.com/anoriar/gophkeeper/internal/server/entry/dto/collection"
	"github.com/anoriar/gophkeeper/internal/server/entry/dto/response/sync"
	"github.com/anoriar/gophkeeper/internal/server/entry/entity"
)

type SyncResponseFactory struct {
}

func NewSyncResponseFactory() *SyncResponseFactory {
	return &SyncResponseFactory{}
}

func (f *SyncResponseFactory) CreateSyncResponse(entryCollection collection.EntryCollection, syncType enum.EntryType) sync.SyncResponse {
	syncResponseItems := make([]sync.SyncResponseItem, 0, len(entryCollection.Entries))
	for _, entry := range entryCollection.Entries {
		responseItem := f.CreateSyncResponseItem(entry)
		syncResponseItems = append(syncResponseItems, responseItem)
	}
	return *sync.NewSyncResponse(syncResponseItems, syncType)
}

func (f *SyncResponseFactory) CreateSyncResponseItem(entry entity.Entry) sync.SyncResponseItem {
	return *sync.NewSyncResponseItem(entry.OriginalId, entry.UpdatedAt, base64.StdEncoding.EncodeToString(entry.Data), entry.Meta)
}
