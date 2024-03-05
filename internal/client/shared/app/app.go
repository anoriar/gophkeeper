package app

import (
	"github.com/anoriar/gophkeeper/internal/client/user/services/auth"
	"go.uber.org/zap"

	"github.com/anoriar/gophkeeper/internal/client/shared/config"
	loggerPkg "github.com/anoriar/gophkeeper/internal/server/shared/app/logger"
)

// App missing godoc.
type App struct {
	Config      *config.Config
	Logger      *zap.Logger
	AuthService auth.AuthServiceInterface
}

// NewApp missing godoc.
func NewApp(cnf *config.Config) (*App, error) {

	logger, err := loggerPkg.Initialize(cnf.LogLevel)
	if err != nil {
		return nil, err
	}

	authService := auth.NewAuthService()

	return &App{
		Config:      cnf,
		Logger:      logger,
		AuthService: authService,
	}, nil
}

func (app *App) Close() {
	app.Logger.Sync()
}
