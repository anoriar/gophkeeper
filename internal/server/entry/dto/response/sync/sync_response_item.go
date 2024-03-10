package sync

import (
	"encoding/json"
	"time"

	"github.com/anoriar/gophkeeper/internal/server/entry/enum"
)

type SyncResponseItem struct {
	OriginalId string          `json:"originalId"`
	EntryType  enum.EntryType  `json:"type"`
	UpdatedAt  time.Time       `json:"updatedAt"`
	Data       string          `json:"data"`
	Meta       json.RawMessage `json:"meta"`
}

func NewSyncResponseItem(originalId string, entryType enum.EntryType, updatedAt time.Time, data string, meta json.RawMessage) *SyncResponseItem {
	return &SyncResponseItem{OriginalId: originalId, EntryType: entryType, UpdatedAt: updatedAt, Data: data, Meta: meta}
}
