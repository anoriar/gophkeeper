package entry_ext

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
	sharedErrors "github.com/anoriar/gophkeeper/internal/client/shared/errors"
)

type SyncResponseItem struct {
	// OriginalId - id записи на клиенте
	OriginalId string
	// EntryType - тип
	EntryType enum.EntryType
	// UpdatedAt - время обновления
	UpdatedAt time.Time
	// Data - зашифрованные данные в base64
	Data string
	// Meta - метаданные
	Meta json.RawMessage
}

func (s *SyncResponseItem) UnmarshalJSON(data []byte) error {
	var alias struct {
		OriginalId string          `json:"originalId"`
		EntryType  string          `json:"type"`
		UpdatedAt  string          `json:"updatedAt"`
		Data       string          `json:"data"`
		Meta       json.RawMessage `json:"meta"`
	}
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	s.OriginalId = alias.OriginalId

	if !enum.IsEntryType(alias.EntryType) {
		return fmt.Errorf("%w: invalid EntryType value: %s", sharedErrors.ErrInternalError, alias.EntryType)
	}
	s.EntryType = (enum.EntryType)(alias.EntryType)
	updatedAt, err := time.Parse(time.RFC3339, alias.UpdatedAt)
	if err != nil {
		return err
	}
	s.UpdatedAt = updatedAt
	s.Meta = alias.Meta

	s.Data = alias.Data

	return nil
}
