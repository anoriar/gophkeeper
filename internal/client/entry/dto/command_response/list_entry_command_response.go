package command_response

import (
	"time"

	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
)

type ListEntryCommandResponse struct {
	Id        string         `json:"id"`
	EntryType enum.EntryType `json:"type"`
	UpdatedAt time.Time      `json:"updatedAt"`
	IsDeleted bool           `json:"isDeleted"`
}
