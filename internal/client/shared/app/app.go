package app

import (
	"go.uber.org/zap"

	"github.com/anoriar/gophkeeper/internal/client/shared/services/uuid"

	entryFactoryPkg "github.com/anoriar/gophkeeper/internal/client/entry/factory"
	"github.com/anoriar/gophkeeper/internal/client/entry/repository/entry_ext"

	"github.com/anoriar/gophkeeper/internal/client/entry/services/encoder"

	entryRepositoryPkg "github.com/anoriar/gophkeeper/internal/client/entry/repository/entry"
	"github.com/anoriar/gophkeeper/internal/client/entry/services/entry"
	"github.com/anoriar/gophkeeper/internal/client/entry/services/service_provider"

	"github.com/anoriar/gophkeeper/internal/client/shared/app/client"
	"github.com/anoriar/gophkeeper/internal/client/user/repository/secret"
	"github.com/anoriar/gophkeeper/internal/client/user/repository/user"
	"github.com/anoriar/gophkeeper/internal/client/user/services/auth"

	loggerPkg "github.com/anoriar/gophkeeper/internal/client/shared/app/logger"
	"github.com/anoriar/gophkeeper/internal/client/shared/config"
)

const (
	loginFile = "./.data/entries/logins"
	cardFile  = "./.data/entries/cards"
	textFile  = "./.data/entries/texts"
	binFile   = "./.data/entries/binaries"
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

	uuidGen := uuid.NewUUIDGenerator()
	gophkeeperHttpClient := client.NewHTTPClient(cnf.ServerAddress, logger)

	userRepository := user.NewUserRepository(gophkeeperHttpClient)
	secretRepository := secret.NewSecretRepository()
	authService := auth.NewAuthService(userRepository, secretRepository, logger)

	aesEncoder := encoder.NewAesDataEncoder()

	loginEntryRepository := entryRepositoryPkg.NewEntrySingleFileRepository(loginFile)
	cardEntryRepository := entryRepositoryPkg.NewEntrySingleFileRepository(cardFile)
	textEntryRepository := entryRepositoryPkg.NewEntrySingleFileRepository(textFile)
	binEntryRepository := entryRepositoryPkg.NewEntrySingleFileRepository(binFile)

	extEntryRepository := entry_ext.NewEntryExtRepository(gophkeeperHttpClient)

	loginEntryService := entry.NewEntryService(
		entryFactoryPkg.NewEntryFactory(uuidGen),
		loginEntryRepository,
		secretRepository,
		aesEncoder,
		extEntryRepository,
		logger,
	)
	cardEntryService := entry.NewEntryService(
		entryFactoryPkg.NewEntryFactory(uuidGen),
		cardEntryRepository,
		secretRepository,
		aesEncoder,
		extEntryRepository,
		logger,
	)
	textEntryService := entry.NewEntryService(
		entryFactoryPkg.NewEntryFactory(uuidGen),
		textEntryRepository,
		secretRepository,
		aesEncoder,
		extEntryRepository,
		logger,
	)

	binEntryService := entry.NewEntryService(
		entryFactoryPkg.NewEntryFactory(uuidGen),
		binEntryRepository,
		secretRepository,
		aesEncoder,
		extEntryRepository,
		logger,
	)

	entryServiceProvider := service_provider.NewEntryServiceProvider(loginEntryService, cardEntryService, textEntryService, binEntryService)

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
