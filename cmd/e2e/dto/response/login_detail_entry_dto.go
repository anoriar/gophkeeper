package response

import (
	"encoding/json"
	"time"
)

type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginDetailEntryResponse struct {
	Id        string          `json:"id"`
	EntryType string          `json:"type"`
	UpdatedAt time.Time       `json:"updatedAt"`
	IsDeleted bool            `json:"isDeleted"`
	Data      LoginData       `json:"data"`
	Meta      json.RawMessage `json:"meta"`
}
