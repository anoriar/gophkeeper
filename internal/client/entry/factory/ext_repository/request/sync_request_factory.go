package request

import (
	"github.com/anoriar/gophkeeper/internal/client/entry/dto/repository/entry_ext"
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
)

type SyncRequestFactory struct {
}

func NewSyncRequestFactory() *SyncRequestFactory {
	return &SyncRequestFactory{}
}

func (f *SyncRequestFactory) CreateFromEntries(entries []entity.Entry) []entry_ext.SyncRequestItem {
	requestItems := make([]entry_ext.SyncRequestItem, 0, len(entries))
	for _, entryEntity := range entries {
		requestItems = append(requestItems, entry_ext.SyncRequestItem{
			OriginalId: entryEntity.Id,
			EntryType:  entryEntity.EntryType,
			UpdatedAt:  entryEntity.UpdatedAt,
			IsDeleted:  entryEntity.IsDeleted,
			Data:       entryEntity.Data,
			Meta:       entryEntity.Meta,
		})
	}
	return requestItems
}
