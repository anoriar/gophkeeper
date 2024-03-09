package factory

import (
	"encoding/json"
	"github.com/google/uuid"

	"github.com/anoriar/gophkeeper/internal/server/entry/dto/request/sync"
	"github.com/anoriar/gophkeeper/internal/server/entry/entity"
)

type EntryFactory struct {
}

func NewEntryFactory() *EntryFactory {
	return &EntryFactory{}
}

func (f *EntryFactory) CreateNewEntryFromRequestItem(requestItem sync.SyncRequestItem, userID string) (entity.Entry, error) {
	data, err := json.Marshal(requestItem.Data)
	if err != nil {
		return entity.Entry{}, err
	}
	return *entity.NewEntry(
		uuid.NewString(),
		requestItem.OriginalId,
		userID,
		requestItem.EntryType,
		requestItem.UpdatedAt,
		data,
		requestItem.Meta,
	), nil
}

func (f *EntryFactory) CreateEntryFromRequestItem(id string, requestItem sync.SyncRequestItem, userID string) (entity.Entry, error) {
	data, err := json.Marshal(requestItem.Data)
	if err != nil {
		return entity.Entry{}, err
	}
	return *entity.NewEntry(
		id,
		requestItem.OriginalId,
		userID,
		requestItem.EntryType,
		requestItem.UpdatedAt,
		data,
		requestItem.Meta,
	), nil
}
