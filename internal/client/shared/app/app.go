package app

import (
	entryFactoryPkg "github.com/anoriar/gophkeeper/internal/client/entry/factory/entry"
	entryRepositoryPkg "github.com/anoriar/gophkeeper/internal/client/entry/repository/entry"
	"github.com/anoriar/gophkeeper/internal/client/entry/services/entry"
	"github.com/anoriar/gophkeeper/internal/client/entry/services/service_provider"
	"go.uber.org/zap"

	"github.com/anoriar/gophkeeper/internal/client/shared/app/client"
	"github.com/anoriar/gophkeeper/internal/client/user/repository/secret"
	"github.com/anoriar/gophkeeper/internal/client/user/repository/user"
	"github.com/anoriar/gophkeeper/internal/client/user/services/auth"

	loggerPkg "github.com/anoriar/gophkeeper/internal/client/shared/app/logger"
	"github.com/anoriar/gophkeeper/internal/client/shared/config"
)

// App missing godoc.
type App struct {
	Config               *config.Config
	Logger               *zap.Logger
	AuthService          auth.AuthServiceInterface
	EntryServiceProvider service_provider.EntryServiceProviderInterface
}

// NewApp missing godoc.
func NewApp(cnf *config.Config) (*App, error) {

	logger, err := loggerPkg.Initialize(cnf.LogLevel)
	if err != nil {
		return nil, err
	}

	gophkeeperHttpClient := client.NewHTTPClient(cnf.ServerAddress, logger)

	userRepository := user.NewUserRepository(gophkeeperHttpClient)
	secretRepository := secret.NewSecretRepository()
	authService := auth.NewAuthService(userRepository, secretRepository, logger)

	loginEntryRepository := entryRepositoryPkg.NewEntrySingleFileRepository("login")
	cardEntryRepository := entryRepositoryPkg.NewEntrySingleFileRepository("card")

	loginEntryService := entry.NewLoginEntryService(entryFactoryPkg.NewEntryFactory(), loginEntryRepository)
	//TODO: проставить другой сервис
	cardEntryService := entry.NewLoginEntryService(entryFactoryPkg.NewEntryFactory(), cardEntryRepository)

	entryServiceProvider := service_provider.NewEntryServiceProvider(loginEntryService, cardEntryService)

	return &App{
		Config:               cnf,
		Logger:               logger,
		AuthService:          authService,
		EntryServiceProvider: entryServiceProvider,
	}, nil
}

func (app *App) Close() {
	app.Logger.Sync()
}
