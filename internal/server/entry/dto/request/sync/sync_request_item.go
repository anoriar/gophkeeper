package sync

import (
	"encoding/json"
	"fmt"
	"time"

	errors2 "github.com/anoriar/gophkeeper/internal/server/entry/errors"

	"github.com/anoriar/gophkeeper/internal/server/entry/enum"
)

type SyncRequestItem struct {
	OriginalId string
	EntryType  enum.EntryType
	UpdatedAt  time.Time
	Data       string
	Meta       json.RawMessage
	IsDeleted  bool
}

func (e *SyncRequestItem) UnmarshalJSON(data []byte) error {
	var alias struct {
		OriginalId string          `json:"originalId"`
		EntryType  string          `json:"type"`
		UpdatedAt  string          `json:"updatedAt"`
		Data       string          `json:"data"`
		Meta       json.RawMessage `json:"meta"`
		IsDeleted  bool            `json:"isDeleted"`
	}
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	e.OriginalId = alias.OriginalId
	e.IsDeleted = alias.IsDeleted

	if !enum.IsEntryType(alias.EntryType) {
		return fmt.Errorf("%w: invalid EntryType value: %s", errors2.ErrSyncRequestNotValid, alias.EntryType)
	}
	e.EntryType = (enum.EntryType)(alias.EntryType)
	updatedAt, err := time.Parse(time.RFC3339, alias.UpdatedAt)
	if err != nil {
		return err
	}
	e.UpdatedAt = updatedAt
	e.Meta = alias.Meta

	e.Data = alias.Data

	return nil
}
