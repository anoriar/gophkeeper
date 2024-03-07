package main

import (
	"context"
	"fmt"
	"log"

	appPkg "github.com/anoriar/gophkeeper/internal/client/shared/app"
	"github.com/anoriar/gophkeeper/internal/client/shared/config"
	commandPkg "github.com/anoriar/gophkeeper/internal/client/shared/services/command"
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
	defer app.Close()

	command, err := ParseFlags()
	if err != nil {
		log.Fatalf("parse command error %v", err.Error())
	}
	errs := command.Validate()
	if len(errs) > 0 {
		log.Fatalf("validation error %v", err.Error())
	}

	cmdExecutor := commandPkg.NewCommandExecutor(app)
	err = cmdExecutor.ExecuteCommand(context.Background(), command)
	if err != nil {
		fmt.Printf("status: failed\n%v\n", err.Error())
	} else {
		fmt.Printf("status: success")
	}
}
