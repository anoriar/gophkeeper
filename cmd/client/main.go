package main

import (
	appPkg "github.com/anoriar/gophkeeper/internal/client/shared/app"
	"github.com/anoriar/gophkeeper/internal/client/shared/config"
	commandPkg "github.com/anoriar/gophkeeper/internal/client/shared/services/command"
	"log"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("recovering from panic: %v", r)
		}
	}()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config error %v", err.Error())
	}

	app, err := appPkg.NewApp(cfg)
	if err != nil {
		log.Fatalf("init app error %v", err.Error())
	}

	command, err := ParseFlags()
	if err != nil {
		log.Fatalf("parse command error %v", err.Error())
	}
	errs := command.Validate()
	if len(errs) > 0 {
		log.Fatalf("validation error %v", err.Error())
	}

	cmdExecutor := commandPkg.NewCommandExecutor(app)
	err = cmdExecutor.ExecuteCommand(command)
	if err != nil {
		log.Fatalf("execute command error %v", err.Error())
	}

}
