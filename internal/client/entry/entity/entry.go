package entity

import (
	"encoding/json"
	"time"

	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
)

type Entry struct {
	Id        string
	EntryType enum.EntryType
	UpdatedAt time.Time
	IsDeleted bool
	Data      []byte
	Meta      json.RawMessage
}

func NewEntry(id string, entryType enum.EntryType, updatedAt time.Time, isDeleted bool, data []byte, meta json.RawMessage) *Entry {
	return &Entry{Id: id, EntryType: entryType, UpdatedAt: updatedAt, IsDeleted: isDeleted, Data: data, Meta: meta}
}
