package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/anoriar/gophkeeper/internal/client/entry/dto"
	entryCommands "github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
	"github.com/anoriar/gophkeeper/internal/client/shared/dto/command"
	userCommands "github.com/anoriar/gophkeeper/internal/client/user/dto/command"
	pflag "github.com/spf13/pflag"
	"os"
)

func ParseFlags() (command.CommandInterface, error) {
	registerFlags := pflag.NewFlagSet("register", pflag.ExitOnError)
	loginFlags := pflag.NewFlagSet("login", pflag.ExitOnError)
	addFlags := pflag.NewFlagSet("add", pflag.ExitOnError)
	editFlags := pflag.NewFlagSet("edit", pflag.ExitOnError)
	deleteFlags := pflag.NewFlagSet("delete", pflag.ExitOnError)
	listFlags := pflag.NewFlagSet("list", pflag.ExitOnError)
	detailFlags := pflag.NewFlagSet("detail", pflag.ExitOnError)
	syncFlags := pflag.NewFlagSet("sync", pflag.ExitOnError)

	if len(os.Args) <= 1 {
		exitWithError(fmt.Errorf("not valid command"))
	}

	switch os.Args[1] {
	case "register":
		reg, err := parseRegisterCommand(registerFlags)
		if err != nil {
			return nil, fmt.Errorf("register command: %v", err)
		}
		return reg, nil
	case "login":
		login, err := parseLoginCommand(loginFlags)
		if err != nil {
			return nil, fmt.Errorf("login command: %v", err)
		}
		return login, nil
	case "add":
		add, err := parseAddEntryCommand(addFlags)
		if err != nil {
			return nil, fmt.Errorf("add command: %v", err)
		}
		return add, err
	case "edit":
		edit, err := parseEditEntryCommand(editFlags)
		if err != nil {
			return nil, fmt.Errorf("edit command: %v", err)
		}
		return edit, err
	case "delete":
		deleteCommand, err := parseDeleteEntryCommand(deleteFlags)
		if err != nil {
			return nil, fmt.Errorf("delete command: %v", err)
		}
		return deleteCommand, err
	case "list":
		listCommand, err := parseListEntryCommand(listFlags)
		if err != nil {
			return nil, fmt.Errorf("list command: %v", err)
		}
		return listCommand, err
	case "detail":
		detailCommand, err := parseDetailEntryCommand(detailFlags)
		if err != nil {
			return nil, fmt.Errorf("detail command: %v", err)
		}
		return detailCommand, err
	case "sync":
		syncCommand, err := parseSyncEntryCommand(syncFlags)
		if err != nil {
			return nil, fmt.Errorf("sync command: %v", err)
		}
		return syncCommand, nil
	default:
		return nil, fmt.Errorf("not valid command")
	}
}

func exitWithError(err error) {
	fmt.Printf("Error: %s\n", err.Error())
	flag.Usage()
	os.Exit(1)
}

func parseRegisterCommand(flags *pflag.FlagSet) (*userCommands.RegisterCommand, error) {
	registerCommand := &userCommands.RegisterCommand{}
	flags.StringVarP(&registerCommand.UserName, "user", "u", "", "Username")
	flags.StringVarP(&registerCommand.Password, "pass", "p", "", "Password")
	flags.StringVarP(&registerCommand.MasterPassword, "masterpass", "m", "", "Master password")

	err := flags.Parse(os.Args[2:])
	if err != nil {
		return nil, err
	}
	errs := registerCommand.Validate()
	if errs != nil {
		return nil, fmt.Errorf("validation error:\n%s", errs.String())
	}
	return registerCommand, nil
}

func parseLoginCommand(flags *pflag.FlagSet) (*userCommands.LoginCommand, error) {
	loginCommand := &userCommands.LoginCommand{}
	flags.StringVarP(&loginCommand.UserName, "user", "u", "", "Username")
	flags.StringVarP(&loginCommand.Password, "pass", "p", "", "Password")
	flags.StringVarP(&loginCommand.MasterPassword, "masterpass", "m", "", "Master password")

	err := flags.Parse(os.Args[2:])
	if err != nil {
		return nil, err
	}
	errs := loginCommand.Validate()
	if errs != nil {
		return nil, fmt.Errorf("validation error:\n%s", errs.String())
	}
	return loginCommand, nil
}

func parseAddEntryCommand(flags *pflag.FlagSet) (*entryCommands.AddEntryCommand, error) {
	var entryTypeStr string
	var dataStr string
	var metaStr string

	flags.StringVar(&entryTypeStr, "t", "", "type")
	flags.StringVar(&dataStr, "d", "", "data")
	flags.StringVar(&metaStr, "m", "", "meta")
	err := flags.Parse(os.Args[2:])
	if err != nil {
		return nil, err
	}

	entryType, data, meta, err := parseDataAndEntryType(entryTypeStr, dataStr, metaStr)
	if err != nil {
		return nil, err
	}

	entryCommand := &entryCommands.AddEntryCommand{}

	entryCommand.EntryType = entryType
	entryCommand.Data = data
	entryCommand.Meta = meta

	return entryCommand, nil
}

func parseEditEntryCommand(flags *pflag.FlagSet) (*entryCommands.EditEntryCommand, error) {
	var id string
	var entryTypeStr string
	var dataStr string
	var metaStr string

	flags.StringVar(&id, "i", "", "id")
	flags.StringVar(&entryTypeStr, "t", "", "type")
	flags.StringVar(&dataStr, "d", "", "data")
	flags.StringVar(&metaStr, "m", "", "meta")
	err := flags.Parse(os.Args[2:])
	if err != nil {
		return nil, err
	}

	entryType, data, meta, err := parseDataAndEntryType(entryTypeStr, dataStr, metaStr)
	if err != nil {
		return nil, err
	}

	entryCommand := &entryCommands.EditEntryCommand{}

	entryCommand.Id = id
	entryCommand.EntryType = entryType
	entryCommand.Data = data
	entryCommand.Meta = meta

	return entryCommand, nil
}

func parseEntryType(entryType string) (enum.EntryType, error) {
	switch entryType {
	case string(enum.Login):
		return enum.Login, nil
	case string(enum.Card):
		return enum.Card, nil
	default:
		return "", errors.New("not valid entry type")
	}
}

func parseDataAndEntryType(entryType string, data string, metaStr string) (enum.EntryType, interface{}, json.RawMessage, error) {
	var meta json.RawMessage
	if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
		return "", nil, json.RawMessage{}, err
	}

	switch entryType {
	case string(enum.Login):
		var loginData dto.LoginData
		if err := json.Unmarshal([]byte(data), &loginData); err != nil {
			return "", nil, json.RawMessage{}, err
		}
		return enum.Login, loginData, meta, nil
	case string(enum.Card):
		var cardData dto.CardData
		if err := json.Unmarshal([]byte(data), &cardData); err != nil {
			return "", nil, json.RawMessage{}, err
		}
		return enum.Card, cardData, meta, nil
	default:
		return "", nil, json.RawMessage{}, errors.New("not valid entry type")
	}
}

func parseListEntryCommand(flags *pflag.FlagSet) (*entryCommands.ListEntryCommand, error) {
	var entryTypeStr string

	flags.StringVar(&entryTypeStr, "t", "", "type")
	err := flags.Parse(os.Args[2:])
	if err != nil {
		return nil, err
	}

	entryType, err := parseEntryType(entryTypeStr)
	if err != nil {
		return nil, err
	}

	entryCommand := &entryCommands.ListEntryCommand{}
	entryCommand.EntryType = entryType

	return entryCommand, nil
}

func parseDeleteEntryCommand(flags *pflag.FlagSet) (*entryCommands.DeleteEntryCommand, error) {
	var id string
	var entryTypeStr string

	flags.StringVar(&id, "i", "", "id")
	flags.StringVar(&entryTypeStr, "t", "", "type")
	err := flags.Parse(os.Args[2:])
	if err != nil {
		return nil, err
	}

	entryType, err := parseEntryType(entryTypeStr)
	if err != nil {
		return nil, err
	}

	entryCommand := &entryCommands.DeleteEntryCommand{}

	entryCommand.Id = id
	entryCommand.EntryType = entryType

	return entryCommand, nil
}

func parseDetailEntryCommand(flags *pflag.FlagSet) (*entryCommands.DetailEntryCommand, error) {
	var id string
	var entryTypeStr string

	flags.StringVar(&id, "i", "", "id")
	flags.StringVar(&entryTypeStr, "t", "", "type")
	err := flags.Parse(os.Args[2:])
	if err != nil {
		return nil, err
	}

	entryType, err := parseEntryType(entryTypeStr)
	if err != nil {
		return nil, err
	}

	entryCommand := &entryCommands.DetailEntryCommand{}

	entryCommand.Id = id
	entryCommand.EntryType = entryType

	return entryCommand, nil
}

func parseSyncEntryCommand(flags *pflag.FlagSet) (*entryCommands.SyncEntryCommand, error) {
	entryCommand := &entryCommands.SyncEntryCommand{}

	return entryCommand, nil
}
