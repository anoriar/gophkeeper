package factory

import (
	"encoding/base64"
	uuid2 "github.com/anoriar/gophkeeper/internal/server/shared/services/uuid"

	"github.com/anoriar/gophkeeper/internal/server/entry/dto/request/sync"
	"github.com/anoriar/gophkeeper/internal/server/entry/entity"
)

type EntryFactory struct {
	uuidGen uuid2.UUIDGeneratorInterface
}

func NewEntryFactory(uuidGen uuid2.UUIDGeneratorInterface) *EntryFactory {
	return &EntryFactory{uuidGen: uuidGen}
}

func (f *EntryFactory) CreateNewEntryFromRequestItem(requestItem sync.SyncRequestItem, userID string) (entity.Entry, error) {
	data, err := base64.StdEncoding.DecodeString(requestItem.Data)
	if err != nil {
		return entity.Entry{}, err
	}
	return *entity.NewEntry(
		f.uuidGen.NewString(),
		requestItem.OriginalId,
		userID,
		requestItem.EntryType,
		requestItem.UpdatedAt,
		data,
		requestItem.Meta,
	), nil
}

func (f *EntryFactory) CreateEntryFromRequestItem(id string, requestItem sync.SyncRequestItem, userID string) (entity.Entry, error) {
	data, err := base64.StdEncoding.DecodeString(requestItem.Data)
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
