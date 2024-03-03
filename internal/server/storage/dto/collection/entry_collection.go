package collection

import "github.com/anoriar/gophkeeper/internal/server/storage/entity"

type EntryCollection struct {
	entries []entity.Entry
}

func NewEntryCollection(entries []entity.Entry) *EntryCollection {
	return &EntryCollection{
		entries: entries,
	}
}

func (c *EntryCollection) Contains(id string) bool {
	for _, entry := range c.entries {
		if entry.Id == id {
			return true
		}
	}
	return false
}

func (c *EntryCollection) FindById(id string) *entity.Entry {
	for _, entry := range c.entries {
		if entry.Id == id {
			return &entry
		}
	}
	return nil
}
