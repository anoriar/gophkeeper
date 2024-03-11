package entry_ext

import (
	"encoding/json"
	"time"
)

type SyncResponseItem struct {
	// OriginalId - id записи на клиенте
	OriginalId string
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
		UpdatedAt  string          `json:"updatedAt"`
		Data       string          `json:"data"`
		Meta       json.RawMessage `json:"meta"`
	}
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	s.OriginalId = alias.OriginalId

	updatedAt, err := time.Parse(time.RFC3339, alias.UpdatedAt)
	if err != nil {
		return err
	}
	s.UpdatedAt = updatedAt
	s.Meta = alias.Meta

	s.Data = alias.Data

	return nil
}
