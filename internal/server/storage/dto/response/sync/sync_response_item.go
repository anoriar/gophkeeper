package sync

import (
	"encoding/json"
	"time"

	"github.com/anoriar/gophkeeper/internal/server/storage/enum"
)

type SyncResponseItem struct {
	Id        string
	EntryType enum.EntryType
	UpdatedAt time.Time
	Data      interface{}
	Meta      json.RawMessage
}

func NewSyncResponseItem(id string, entryType enum.EntryType, updatedAt time.Time, data interface{}, meta json.RawMessage) *SyncResponseItem {
	return &SyncResponseItem{Id: id, EntryType: entryType, UpdatedAt: updatedAt, Data: data, Meta: meta}
}
