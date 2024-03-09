package command

import (
	"context"
	"encoding/json"
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
		return nil
	case *userCommandPkg.LoginCommand:
		if cmd, ok := command.(*userCommandPkg.LoginCommand); ok {
			return sp.app.AuthService.Login(ctx, *cmd)
		}
		return nil
	case *entryCommandPkg.AddEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.AddEntryCommand); ok {
			return sp.app.EntryServiceProvider.Add(ctx, *cmd)
		}
		return nil
	case *entryCommandPkg.EditEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.EditEntryCommand); ok {
			return sp.app.EntryServiceProvider.Edit(ctx, *cmd)
		}
		return nil
	case *entryCommandPkg.DeleteEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.DeleteEntryCommand); ok {
			return sp.app.EntryServiceProvider.Delete(ctx, *cmd)
		}
		return nil
	case *entryCommandPkg.DetailEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.DetailEntryCommand); ok {
			entry, err := sp.app.EntryServiceProvider.Detail(ctx, *cmd)
			if err != nil {
				return err
			}

			response, err := json.MarshalIndent(entry, "", "    ")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", response)
		}
		return nil
	case *entryCommandPkg.ListEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.ListEntryCommand); ok {
			entries, err := sp.app.EntryServiceProvider.GetList(ctx, *cmd)
			if err != nil {
				return err
			}

			response, err := json.MarshalIndent(entries, "", "    ")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", response)
		}
		return nil
	case *entryCommandPkg.SyncEntryCommand:
		if cmd, ok := command.(*entryCommandPkg.SyncEntryCommand); ok {
			err := sp.app.EntryServiceProvider.Sync(ctx, *cmd)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("command not exists")
	}
	return fmt.Errorf("command not exists")
}
