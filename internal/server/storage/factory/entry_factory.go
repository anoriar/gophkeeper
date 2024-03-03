package factory

import (
	"encoding/json"
	"github.com/anoriar/gophkeeper/internal/server/storage/dto/request"
	"github.com/anoriar/gophkeeper/internal/server/storage/entity"
)

type EntryFactory struct {
}

func NewEntryFactory() *EntryFactory {
	return &EntryFactory{}
}

func (f *EntryFactory) CreateEntryFromRequestItem(requestItem request.SyncRequestItem, userID string) (entity.Entry, error) {
	data, err := json.Marshal(requestItem.Data)
	if err != nil {
		return entity.Entry{}, err
	}
	return *entity.NewEntry(
		requestItem.Id,
		userID,
		requestItem.EntryType,
		requestItem.UpdatedAt,
		data,
		requestItem.Meta,
	), nil
}
