package main

import (
	"log"

	appPkg "github.com/anoriar/gophkeeper/internal/server/shared/app"
	"github.com/anoriar/gophkeeper/internal/server/shared/config"
	"github.com/anoriar/gophkeeper/internal/server/shared/router"
	"github.com/anoriar/gophkeeper/internal/server/shared/server"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("create config error %v", err.Error())
	}

	app, err := appPkg.NewApp(conf)
	if err != nil {
		log.Fatalf("init app error %v", err.Error())
	}
	defer app.Close()

	r := router.NewRouter(app)

	err = server.RunServer(app, r)
	if err != nil {
		log.Fatalf("init router error %v", err.Error())
	}
}
