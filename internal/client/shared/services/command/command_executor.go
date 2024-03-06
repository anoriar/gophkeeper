package command

import (
	"fmt"

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
	default:
		return fmt.Errorf("command not exists")
	}
	return fmt.Errorf("command not exists")
}
