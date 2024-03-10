package factory

import (
	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/dto/repository/entry_ext"
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
)

type EntryFactoryInterface interface {
	CreateFromAddCmd(command command.AddEntryCommand) (entity.Entry, error)
	CreateFromEditCmd(command command.EditEntryCommand) (entity.Entry, error)
	CreateFromSyncResponse(syncResponseItems []entry_ext.SyncResponseItem) ([]entity.Entry, error)
}
