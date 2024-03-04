package app

import (
	"go.uber.org/zap"

	entryRepositoryPkg "github.com/anoriar/gophkeeper/internal/server/storage/repository"
	"github.com/anoriar/gophkeeper/internal/server/storage/services/sync"

	dbPkg "github.com/anoriar/gophkeeper/internal/server/shared/app/db"
	loggerPkg "github.com/anoriar/gophkeeper/internal/server/shared/app/logger"
	"github.com/anoriar/gophkeeper/internal/server/shared/config"
	userRepositoryPkg "github.com/anoriar/gophkeeper/internal/server/user/repository"
	"github.com/anoriar/gophkeeper/internal/server/user/services/auth"
)

// App missing godoc.
type App struct {
	Config      *config.Config
	Logger      *zap.Logger
	Database    dbPkg.DatabaseInterface
	AuthService auth.AuthServiceInterface
	SyncService sync.SyncServiceInterface
}

// NewApp missing godoc.
func NewApp(cnf *config.Config) (*App, error) {

	db, err := dbPkg.InitializeDatabase(cnf.DatabaseURI)
	if err != nil {
		return nil, err
	}
	logger, err := loggerPkg.Initialize(cnf.LogLevel)
	if err != nil {
		return nil, err
	}

	userRepository := userRepositoryPkg.NewUserRepository(db)
	authService := auth.NewAuthService(userRepository,
		cnf,
		logger,
	)

	entryRepository := entryRepositoryPkg.NewEntryRepository(db)
	syncService := sync.NewSyncService(entryRepository, db, logger)

	return &App{
		Config:      cnf,
		Logger:      logger,
		Database:    db,
		AuthService: authService,
		SyncService: syncService,
	}, nil
}

func (app *App) Close() {
	app.Database.Close()
	app.Logger.Sync()
}
