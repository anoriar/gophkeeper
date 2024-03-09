package service_provider

import (
	"context"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command_response"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
)

type EntryServiceProviderInterface interface {
	Add(ctx context.Context, cmd command.AddEntryCommand) error
	Edit(ctx context.Context, cmd command.EditEntryCommand) error
	Delete(ctx context.Context, cmd command.DeleteEntryCommand) error
	Detail(ctx context.Context, cmd command.DetailEntryCommand) (command_response.DetailEntryCommandResponse, error)
	GetList(ctx context.Context, cmd command.ListEntryCommand) ([]command_response.ListEntryCommandResponse, error)
	Sync(ctx context.Context, cmd command.SyncEntryCommand) error
}
