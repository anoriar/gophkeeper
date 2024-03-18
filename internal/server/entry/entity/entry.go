package entity

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/anoriar/gophkeeper/internal/server/entry/enum"
)

type Entry struct {
	Id         string          `db:"id"`
	OriginalId string          `db:"original_id"`
	UserId     string          `db:"user_id"`
	EntryType  enum.EntryType  `db:"type"`
	UpdatedAt  time.Time       `db:"updated_at"`
	Data       []byte          `db:"data"`
	Meta       json.RawMessage `db:"meta"`
}

func NewEntry(id string, originalId string, userId string, entryType enum.EntryType, updatedAt time.Time, data []byte, meta json.RawMessage) *Entry {
	return &Entry{Id: id, OriginalId: originalId, UserId: userId, EntryType: entryType, UpdatedAt: updatedAt, Data: data, Meta: meta}
}

func (e Entry) Equals(other Entry) (bool, error) {
	meta, err := e.Meta.MarshalJSON()
	if err != nil {
		return false, err
	}
	otherMeta, err := other.Meta.MarshalJSON()
	if err != nil {
		return false, err
	}

	return e.Id == other.Id &&
		e.OriginalId == other.OriginalId &&
		e.UserId == other.UserId &&
		e.EntryType == other.EntryType &&
		bytes.Equal(e.Data, other.Data) &&
		e.UpdatedAt == other.UpdatedAt &&
		bytes.Equal(meta, otherMeta), nil

}
