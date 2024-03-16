package response

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto"
	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command_response"
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
	errors2 "github.com/anoriar/gophkeeper/internal/server/shared/errors"
)

type EntryResponseFactory struct {
}

func NewEntryResponseFactory() *EntryResponseFactory {
	return &EntryResponseFactory{}
}

func (f *EntryResponseFactory) CreateDetailResponseFromEntity(entry entity.Entry) (command_response.DetailEntryResponse, error) {
	var data interface{}

	switch entry.EntryType {
	case enum.Login:
		data = &dto.LoginData{}
		err := json.Unmarshal(entry.Data, data)
		if err != nil {
			return command_response.DetailEntryResponse{}, fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
		}
	case enum.Card:
		data = &dto.CardData{}
		err := json.Unmarshal(entry.Data, data)
		if err != nil {
			return command_response.DetailEntryResponse{}, fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
		}
	case enum.Text:
		data = string(entry.Data)
	case enum.Bin:
		data = base64.StdEncoding.EncodeToString(entry.Data)
	default:
		return command_response.DetailEntryResponse{}, fmt.Errorf("%w: %v", errors2.ErrInternalError, "data not compatible with any format")
	}

	return command_response.DetailEntryResponse{
		Id:        entry.Id,
		EntryType: entry.EntryType,
		UpdatedAt: entry.UpdatedAt,
		IsDeleted: entry.IsDeleted,
		Data:      data,
		Meta:      entry.Meta,
	}, nil
}

func (f *EntryResponseFactory) CreateListResponseFromEntity(entries []entity.Entry) []command_response.ListEntryCommandResponse {
	responseEntries := make([]command_response.ListEntryCommandResponse, 0, len(entries))

	for _, entryEntity := range entries {
		responseEntries = append(responseEntries, command_response.ListEntryCommandResponse{
			Id:        entryEntity.Id,
			EntryType: entryEntity.EntryType,
			UpdatedAt: entryEntity.UpdatedAt,
			IsDeleted: entryEntity.IsDeleted,
		})
	}
	return responseEntries
}
