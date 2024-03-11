package request

import (
	"encoding/base64"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/repository/entry_ext"
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
)

type SyncRequestFactory struct {
}

func NewSyncRequestFactory() *SyncRequestFactory {
	return &SyncRequestFactory{}
}

func (f *SyncRequestFactory) CreateSyncRequest(syncType enum.EntryType, entries []entity.Entry) entry_ext.SyncRequest {
	items := f.CreateFromEntries(entries)
	return entry_ext.SyncRequest{
		Items:    items,
		SyncType: syncType,
	}
}

func (f *SyncRequestFactory) CreateFromEntries(entries []entity.Entry) []entry_ext.SyncRequestItem {
	requestItems := make([]entry_ext.SyncRequestItem, 0, len(entries))
	for _, entryEntity := range entries {
		requestItems = append(requestItems, entry_ext.SyncRequestItem{
			OriginalId: entryEntity.Id,
			UpdatedAt:  entryEntity.UpdatedAt,
			IsDeleted:  entryEntity.IsDeleted,
			Data:       base64.StdEncoding.EncodeToString(entryEntity.Data),
			Meta:       entryEntity.Meta,
		})
	}
	return requestItems
}
