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
