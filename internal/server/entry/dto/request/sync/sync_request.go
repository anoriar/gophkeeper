package sync

import (
	"encoding/json"
	"fmt"

	"github.com/anoriar/gophkeeper/internal/server/entry/enum"
	errors2 "github.com/anoriar/gophkeeper/internal/server/entry/errors"
)

type SyncRequest struct {
	Items    []SyncRequestItem `json:"items"`
	SyncType enum.EntryType
	UserID   string
}

func (c *SyncRequest) Contains(id string) bool {
	for _, entry := range c.Items {
		if entry.OriginalId == id {
			return true
		}
	}
	return false
}

func (c *SyncRequest) FindById(id string) *SyncRequestItem {
	for _, entry := range c.Items {
		if entry.OriginalId == id {
			return &entry
		}
	}
	return nil
}

func (c *SyncRequest) UnmarshalJSON(data []byte) error {
	type Alias SyncRequest
	alias := &struct {
		*Alias
		SyncType string `json:"syncType"`
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	if !enum.IsEntryType(alias.SyncType) {
		return fmt.Errorf("%w: invalid SyncType value: %s", errors2.ErrSyncRequestNotValid, alias.SyncType)
	}

	c.SyncType = (enum.EntryType)(alias.SyncType)

	return nil
}
