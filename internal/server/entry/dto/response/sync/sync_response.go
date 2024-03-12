package sync

import (
	"github.com/anoriar/gophkeeper/internal/server/entry/enum"
)

type SyncResponse struct {
	Items    []SyncResponseItem `json:"items"`
	SyncType enum.EntryType     `json:"syncType"`
}

func NewSyncResponse(items []SyncResponseItem, syncType enum.EntryType) *SyncResponse {
	return &SyncResponse{Items: items, SyncType: syncType}
}
