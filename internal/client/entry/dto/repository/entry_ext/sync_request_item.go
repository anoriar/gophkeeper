package entry_ext

import (
	"encoding/json"
	"time"
)

type SyncRequestItem struct {
	// OriginalId - id записи на клиенте
	OriginalId string `json:"originalId"`
	// UpdatedAt - время обновления
	UpdatedAt time.Time `json:"updatedAt"`
	// IsDeleted - признак, что нужно удалить запись
	IsDeleted bool `json:"isDeleted"`
	// Data - зашифрованные данные в base64
	Data string `json:"data"`
	// Meta - метаданные
	Meta json.RawMessage `json:"meta"`
}
