package entry_ext

import "github.com/anoriar/gophkeeper/internal/client/entry/enum"

type SyncRequest struct {
	Items    []SyncRequestItem `json:"items"`
	SyncType enum.EntryType    `json:"syncType"`
}
