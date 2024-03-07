package entry

import (
	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
)

type EntryFactoryInterface interface {
	CreateFromAddCmd(command command.AddEntryCommand) (entity.Entry, error)
	CreateFromEditCmd(command command.EditEntryCommand) (entity.Entry, error)
}