package entry_ext

import (
	"encoding/json"
	"fmt"

	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
	sharedErrors "github.com/anoriar/gophkeeper/internal/client/shared/errors"
)

type SyncResponse struct {
	Items    []SyncResponseItem `json:"items"`
	SyncType enum.EntryType     `json:"syncType"`
}

func (c *SyncResponse) UnmarshalJSON(data []byte) error {
	type Alias SyncResponse
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
		return fmt.Errorf("%w: invalid EntryType value: %s", sharedErrors.ErrInternalError, alias.SyncType)
	}

	c.SyncType = (enum.EntryType)(alias.SyncType)

	return nil
}
