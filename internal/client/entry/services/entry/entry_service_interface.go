package entry

import (
	"context"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
)

type EntryServiceInterface interface {
	Add(ctx context.Context, command command.AddEntryCommand) error
	Edit(ctx context.Context, command command.EditEntryCommand) error
	Detail(ctx context.Context, command command.DetailEntryCommand) (entity.Entry, error)
	Delete(ctx context.Context, command command.DeleteEntryCommand) error
	List(ctx context.Context, command command.ListEntryCommand) ([]entity.Entry, error)
}
