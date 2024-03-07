package entity

import (
	"encoding/json"
	"time"

	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
)

type Entry struct {
	Id        string          `json:"id"`
	EntryType enum.EntryType  `json:"type"`
	UpdatedAt time.Time       `json:"updatedAt"`
	IsDeleted bool            `json:"isDeleted"`
	Data      []byte          `json:"data"`
	Meta      json.RawMessage `json:"meta"`
}

func NewEntry(id string, entryType enum.EntryType, updatedAt time.Time, isDeleted bool, data []byte, meta json.RawMessage) *Entry {
	return &Entry{Id: id, EntryType: entryType, UpdatedAt: updatedAt, IsDeleted: isDeleted, Data: data, Meta: meta}
}
