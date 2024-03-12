package sync

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type SyncRequestItem struct {
	OriginalId string
	UpdatedAt  time.Time
	Data       []byte
	Meta       json.RawMessage
	IsDeleted  bool
}

func (e *SyncRequestItem) UnmarshalJSON(data []byte) error {
	var alias struct {
		OriginalId string          `json:"originalId"`
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

	updatedAt, err := time.Parse(time.RFC3339, alias.UpdatedAt)
	if err != nil {
		return err
	}
	e.UpdatedAt = updatedAt
	e.Meta = alias.Meta

	itemData, err := base64.StdEncoding.DecodeString(alias.Data)
	if err != nil {
		return err
	}
	e.Data = itemData

	return nil
}
