package sync

type SyncResponse struct {
	Items []SyncResponseItem `json:"items"`
}

func NewSyncResponse(items []SyncResponseItem) *SyncResponse {
	return &SyncResponse{Items: items}
}
