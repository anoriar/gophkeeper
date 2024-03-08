package request

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto"
	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
	"github.com/anoriar/gophkeeper/internal/client/shared/errors"
)

type EntryFactory struct {
}

func NewEntryFactory() *EntryFactory {
	return &EntryFactory{}
}

func (l *EntryFactory) CreateFromAddCmd(command command.AddEntryCommand) (entity.Entry, error) {
	data, err := l.createData(command.Data)
	if err != nil {
		return entity.Entry{}, err
	}
	return entity.Entry{
		Id:        uuid.NewString(),
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
