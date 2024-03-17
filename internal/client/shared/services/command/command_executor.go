package command

import (
	"context"
	"errors"

	entryCommandPkg "github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/shared/app"
	sharedCommand "github.com/anoriar/gophkeeper/internal/client/shared/dto/command"
	userCommandPkg "github.com/anoriar/gophkeeper/internal/client/user/dto/command"
)

var ErrNotExecuted = errors.New("command did not executed")
var ErrNotExists = errors.New("command not exists")

type CommandExecutor struct {
	app *app.App
}

func NewCommandExecutor(app *app.App) *CommandExecutor {
	return &CommandExecutor{app: app}
}
func (sp *CommandExecutor) prepareCommandResponse(payload interface{}, error error) sharedCommand.CommandResponse {
	status := "success"
	errorStr := ""
	if error != nil {
		errorStr = error.Error()
		status = "fail"
	}
	return sharedCommand.CommandResponse{
		Status:  status,
		Error:   errorStr,
		Payload: payload,
	}
}

func (sp *CommandExecutor) ExecuteCommand(ctx context.Context, command sharedCommand.CommandInterface) sharedCommand.CommandResponse {
	switch command.(type) {
	case *userCommandPkg.RegisterCommand:
		if cmd, ok := command.(*userCommandPkg.RegisterCommand); ok {
			err := sp.app.AuthService.Register(ctx, *cmd)
			return sp.prepareCommandResponse(nil, err)
		}
		return sp.prepareCommandResponse(nil, ErrNotExecuted)
	case *userCommandPkg.LoginCommand:
		if cmd, ok := command.(*userCommandPkg.LoginCommand); ok {
			err := sp.app.AuthService.Login(ctx, *cmd)
			return sp.prepareCommandResponse(nil, err)
		}
		return sp.prepareCommandResponse(nil, ErrNotExecuted)
	case *entryCommandPkg.AddEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.AddEntryCommand); ok {
			entry, err := sp.app.EntryServiceProvider.Add(ctx, *cmd)
			if err != nil {
				return sp.prepareCommandResponse(nil, err)
			}
			return sp.prepareCommandResponse(entry, err)
		}
		return sp.prepareCommandResponse(nil, ErrNotExecuted)
	case *entryCommandPkg.EditEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.EditEntryCommand); ok {
			entry, err := sp.app.EntryServiceProvider.Edit(ctx, *cmd)
			if err != nil {
				return sp.prepareCommandResponse(nil, err)
			}
			return sp.prepareCommandResponse(entry, err)
		}
		return sp.prepareCommandResponse(nil, ErrNotExecuted)
	case *entryCommandPkg.DeleteEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.DeleteEntryCommand); ok {
			err := sp.app.EntryServiceProvider.Delete(ctx, *cmd)
			return sp.prepareCommandResponse(nil, err)
		}
		return sp.prepareCommandResponse(nil, ErrNotExecuted)
	case *entryCommandPkg.DetailEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.DetailEntryCommand); ok {
			entry, err := sp.app.EntryServiceProvider.Detail(ctx, *cmd)
			if err != nil {
				return sp.prepareCommandResponse(nil, err)
			}
			return sp.prepareCommandResponse(entry, err)
		}
		return sp.prepareCommandResponse(nil, ErrNotExecuted)
	case *entryCommandPkg.ListEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.ListEntryCommand); ok {
			entries, err := sp.app.EntryServiceProvider.GetList(ctx, *cmd)
			if err != nil {
				return sp.prepareCommandResponse(nil, err)
			}

			return sp.prepareCommandResponse(entries, err)
		}
		return sp.prepareCommandResponse(nil, ErrNotExecuted)
	case *entryCommandPkg.SyncEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.SyncEntryCommand); ok {
			err := sp.app.EntryServiceProvider.Sync(ctx, *cmd)
			return sp.prepareCommandResponse(nil, err)
		}
		return sp.prepareCommandResponse(nil, ErrNotExecuted)
	default:
		return sp.prepareCommandResponse(nil, ErrNotExists)
	}
}
