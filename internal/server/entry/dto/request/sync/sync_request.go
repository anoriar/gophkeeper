package sync

type SyncRequest struct {
	Items  []SyncRequestItem `json:"items"`
	UserID string
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
