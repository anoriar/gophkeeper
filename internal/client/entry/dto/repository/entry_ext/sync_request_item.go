package entry_ext

import (
	"encoding/json"
	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
	"time"
)

type SyncRequestItem struct {
	// OriginalId - id записи на клиенте
	OriginalId string `json:"originalId"`
	// EntryType - тип
	EntryType enum.EntryType `json:"type"`
	// UpdatedAt - время обновления
	UpdatedAt time.Time `json:"updatedAt"`
	// IsDeleted - признак, что нужно удалить запись
	IsDeleted bool `json:"isDeleted"`
	// Data - зашифрованные данные
	Data []byte `json:"data"`
	// Meta - метаданные
	Meta json.RawMessage `json:"meta"`
}
