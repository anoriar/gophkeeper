package sync

import (
	"encoding/json"
	"time"
)

type SyncResponseItem struct {
	OriginalId string          `json:"originalId"`
	UpdatedAt  time.Time       `json:"updatedAt"`
	Data       string          `json:"data"`
	Meta       json.RawMessage `json:"meta"`
}

func NewSyncResponseItem(originalId string, updatedAt time.Time, data string, meta json.RawMessage) *SyncResponseItem {
	return &SyncResponseItem{OriginalId: originalId, UpdatedAt: updatedAt, Data: data, Meta: meta}
}
