package command

import (
	"context"
	"fmt"

	entryCommandPkg "github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/shared/app"
	"github.com/anoriar/gophkeeper/internal/client/shared/dto/command"
	userCommandPkg "github.com/anoriar/gophkeeper/internal/client/user/dto/command"
)

type CommandExecutor struct {
	app *app.App
}

func NewCommandExecutor(app *app.App) *CommandExecutor {
	return &CommandExecutor{app: app}
}

func (sp *CommandExecutor) ExecuteCommand(ctx context.Context, command command.CommandInterface) error {
	switch command.(type) {
	case *userCommandPkg.RegisterCommand:
		if cmd, ok := command.(*userCommandPkg.RegisterCommand); ok {
			return sp.app.AuthService.Register(ctx, *cmd)
		}
		break
	case *userCommandPkg.LoginCommand:
		if cmd, ok := command.(*userCommandPkg.LoginCommand); ok {
			return sp.app.AuthService.Login(ctx, *cmd)
		}
		break
	case *entryCommandPkg.AddEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.AddEntryCommand); ok {
			return sp.app.EntryServiceProvider.Add(ctx, *cmd)
		}
		break
	case *entryCommandPkg.EditEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.EditEntryCommand); ok {
			return sp.app.EntryServiceProvider.Edit(ctx, *cmd)
		}
		break
	case *entryCommandPkg.DeleteEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.DeleteEntryCommand); ok {
			return sp.app.EntryServiceProvider.Delete(ctx, *cmd)
		}
		break
	default:
		return fmt.Errorf("command not exists")
	}
	return fmt.Errorf("command not exists")
}
