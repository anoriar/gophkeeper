package command

import (
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

func (sp *CommandExecutor) ExecuteCommand(command command.CommandInterface) error {
	switch command.(type) {
	case *userCommandPkg.RegisterCommand:
		if cmd, ok := command.(*userCommandPkg.RegisterCommand); ok {
			return sp.app.AuthService.Register(*cmd)
		}
		break
	case *userCommandPkg.LoginCommand:
		if cmd, ok := command.(*userCommandPkg.LoginCommand); ok {
			return sp.app.AuthService.Login(*cmd)
		}
		break
	case *entryCommandPkg.AddEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.AddEntryCommand); ok {
			return sp.app.EntryServiceProvider.Add(*cmd)
		}
		break
	case *entryCommandPkg.EditEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.EditEntryCommand); ok {
			return sp.app.EntryServiceProvider.Edit(*cmd)
		}
		break
	default:
		return fmt.Errorf("command not exists")
	}
	return fmt.Errorf("command not exists")
}
