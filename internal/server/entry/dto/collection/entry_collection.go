package collection

import "github.com/anoriar/gophkeeper/internal/server/entry/entity"

type EntryCollection struct {
	Entries []entity.Entry
}

func NewEntryCollection(entries []entity.Entry) *EntryCollection {
	return &EntryCollection{
		Entries: entries,
	}
}

func (c *EntryCollection) Contains(id string) bool {
	for _, entry := range c.Entries {
		if entry.OriginalId == id {
			return true
		}
	}
	return false
}

func (c *EntryCollection) FindByOriginalId(id string) *entity.Entry {
	for _, entry := range c.Entries {
		if entry.OriginalId == id {
			return &entry
		}
	}
	return nil
}
