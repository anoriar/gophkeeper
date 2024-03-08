package command_response

import (
	"encoding/json"
	"time"

	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
)

type DetailEntryCommandResponse struct {
	Id        string          `json:"id"`
	EntryType enum.EntryType  `json:"type"`
	UpdatedAt time.Time       `json:"updatedAt"`
	IsDeleted bool            `json:"isDeleted"`
	Data      interface{}     `json:"data"`
	Meta      json.RawMessage `json:"meta"`
}
