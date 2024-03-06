package entry

import (
	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
)

type EntryFactory struct {
}

func NewEntryFactory() *EntryFactory {
	return &EntryFactory{}
}

func (l EntryFactory) CreateFromAddCmd(command command.AddEntryCommand) entity.Entry {
	//TODO implement me
	panic("implement me")
}

func (l EntryFactory) CreateFromEditCmd(command command.EditEntryCommand) entity.Entry {
	//TODO implement me
	panic("implement me")
}
