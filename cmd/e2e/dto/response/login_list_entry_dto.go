package response

import (
	"time"
)

type LoginListEntryResponse struct {
	Id        string    `json:"id"`
	EntryType string    `json:"type"`
	UpdatedAt time.Time `json:"updatedAt"`
	IsDeleted bool      `json:"isDeleted"`
}
