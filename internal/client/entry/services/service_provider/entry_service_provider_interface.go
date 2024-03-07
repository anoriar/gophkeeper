package service_provider

import (
	"context"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
)

type EntryServiceProviderInterface interface {
	Add(ctx context.Context, cmd command.AddEntryCommand) error
	Edit(ctx context.Context, cmd command.EditEntryCommand) error
	Delete(ctx context.Context, cmd command.DeleteEntryCommand) error
	GetById(ctx context.Context, cmd command.DetailEntryCommand) (entity.Entry, error)
	GetList(ctx context.Context, cmd command.ListEntryCommand) ([]entity.Entry, error)
}
