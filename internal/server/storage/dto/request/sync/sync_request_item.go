package sync

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/anoriar/gophkeeper/internal/server/storage/dto"
	"github.com/anoriar/gophkeeper/internal/server/storage/enum"
)

type SyncRequestItem struct {
	Id        string
	EntryType enum.EntryType
	UpdatedAt time.Time
	Data      interface{}
	Meta      json.RawMessage
	IsDeleted bool
}

func (e *SyncRequestItem) UnmarshalJSON(data []byte) error {
	var alias struct {
		Id        string          `json:"id"`
		EntryType string          `json:"type"`
		UpdatedAt string          `json:"updatedAt"`
		Data      json.RawMessage `json:"data"`
		Meta      json.RawMessage `json:"meta"`
		IsDeleted bool            `json:"isDeleted"`
	}
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	e.Id = alias.Id
	e.IsDeleted = alias.IsDeleted

	switch alias.EntryType {
	case string(enum.Login), string(enum.Card):
		e.EntryType = (enum.EntryType)(alias.EntryType)
	default:
		return fmt.Errorf("invalid EntryType value: %s", alias.EntryType)
	}

	updatedAt, err := time.Parse(time.RFC3339, alias.UpdatedAt)
	if err != nil {
		return err
	}
	e.UpdatedAt = updatedAt
	e.Meta = alias.Meta

	switch e.EntryType {
	case enum.Login:
		e.Data = &dto.LoginData{}
	case enum.Card:
		e.Data = &dto.CardData{}
	}

	if e.Data != nil {
		if err := json.Unmarshal(alias.Data, e.Data); err != nil {
			return fmt.Errorf("value of type (%T): unmarshal: %w", e.Data, err)
		}
	}

	return nil
}
