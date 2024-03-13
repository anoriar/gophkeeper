package factory

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anoriar/gophkeeper/internal/client/shared/services/uuid"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/repository/entry_ext"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto"
	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
	"github.com/anoriar/gophkeeper/internal/client/shared/errors"
)

type EntryFactory struct {
	uuidGen uuid.UUIDGeneratorInterface
}

func NewEntryFactory(uuidGen uuid.UUIDGeneratorInterface) *EntryFactory {
	return &EntryFactory{uuidGen: uuidGen}
}

func (l *EntryFactory) CreateFromAddCmd(command command.AddEntryCommand) (entity.Entry, error) {
	data, err := l.createData(command.Data)
	if err != nil {
		return entity.Entry{}, err
	}
	return entity.Entry{
		Id:        l.uuidGen.NewString(),
		EntryType: command.EntryType,
		UpdatedAt: time.Now(),
		IsDeleted: false,
		Data:      data,
		Meta:      command.Meta,
	}, nil
}

func (l *EntryFactory) CreateFromEditCmd(command command.EditEntryCommand) (entity.Entry, error) {
	data, err := l.createData(command.Data)
	if err != nil {
		return entity.Entry{}, err
	}
	return entity.Entry{
		Id:        command.Id,
		EntryType: command.EntryType,
		UpdatedAt: time.Now(),
		IsDeleted: false,
		Data:      data,
		Meta:      command.Meta,
	}, nil
}

func (l *EntryFactory) createData(data interface{}) ([]byte, error) {
	switch data.(type) {
	case dto.LoginData:
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errors.ErrInternalError, err)
		}
		return dataBytes, nil
	case dto.CardData:
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errors.ErrInternalError, err)
		}
		return dataBytes, nil
	default:
		return nil, fmt.Errorf("%w: %v", errors.ErrInternalError, "data type is not implemented")
	}
}

func (l *EntryFactory) CreateFromSyncResponse(syncResponse entry_ext.SyncResponse) ([]entity.Entry, error) {
	entries := make([]entity.Entry, 0, len(syncResponse.Items))

	for _, responseItem := range syncResponse.Items {
		data, err := base64.StdEncoding.DecodeString(responseItem.Data)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errors.ErrInternalError, "data is not decoded")
		}
		entries = append(entries, entity.Entry{
			Id:        responseItem.OriginalId,
			EntryType: syncResponse.SyncType,
			UpdatedAt: responseItem.UpdatedAt,
			IsDeleted: false,
			Data:      data,
			Meta:      responseItem.Meta,
		})
	}
	return entries, nil
}
