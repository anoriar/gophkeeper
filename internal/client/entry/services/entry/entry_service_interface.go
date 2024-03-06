package entry

import (
	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
)

type EntryServiceInterface interface {
	Add(command command.AddEntryCommand) error
	Edit(command command.EditEntryCommand) error
	Detail(command command.DetailEntryCommand) (entity.Entry, error)
	Delete(command command.DeleteEntryCommand) error
	List(command command.ListEntryCommand) ([]entity.Entry, error)
}
