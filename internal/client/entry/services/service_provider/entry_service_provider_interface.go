package service_provider

import (
	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
)

type EntryServiceProviderInterface interface {
	Add(cmd command.AddEntryCommand) error
	Edit(cmd command.EditEntryCommand) error
	Delete(cmd command.DeleteEntryCommand) error
	GetById(cmd command.DetailEntryCommand) (entity.Entry, error)
	GetList(cmd command.ListEntryCommand) ([]entity.Entry, error)
}
