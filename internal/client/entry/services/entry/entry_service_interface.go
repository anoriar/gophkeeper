package entry

import (
	"context"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command_response"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
)

type EntryServiceInterface interface {
	// Add Добавление записи
	Add(ctx context.Context, command command.AddEntryCommand) error
	// Edit Редактирование записи
	Edit(ctx context.Context, command command.EditEntryCommand) error
	// Detail Полная информация по записи (с расшифрованной информацией)
	Detail(ctx context.Context, command command.DetailEntryCommand) (command_response.DetailEntryCommandResponse, error)
	// Delete Удаление записи, isDeleted=true
	Delete(ctx context.Context, command command.DeleteEntryCommand) error
	// List Список записей (без даты)
	List(ctx context.Context) ([]command_response.ListEntryCommandResponse, error)
	Sync(ctx context.Context, command command.SyncEntryCommand) error
}
