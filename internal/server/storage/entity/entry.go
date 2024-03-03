package entity

import (
	"encoding/json"
	"github.com/anoriar/gophkeeper/internal/server/storage/enum"
	"time"
)

type Entry struct {
	Id        string          `db:"id"`
	UserId    string          `db:"user_id"`
	EntryType enum.EntryType  `db:"type"`
	UpdatedAt time.Time       `db:"updated_at"`
	Data      []byte          `db:"data"`
	Meta      json.RawMessage `db:"meta"`
}

func NewEntry(id string, userId string, entryType enum.EntryType, updatedAt time.Time, data []byte, meta json.RawMessage) *Entry {
	return &Entry{Id: id, UserId: userId, EntryType: entryType, UpdatedAt: updatedAt, Data: data, Meta: meta}
}
