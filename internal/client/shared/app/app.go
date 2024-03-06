package app

import (
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

	gophkeeperHttpClient := client.NewHTTPClient(cnf.ServerAddress, logger)

	userRepository := user.NewUserRepository(gophkeeperHttpClient)
	secretRepository := secret.NewSecretRepository()
	authService := auth.NewAuthService(userRepository, secretRepository, logger)

	return &App{
		Config:      cnf,
		Logger:      logger,
		AuthService: authService,
	}, nil
}

func (app *App) Close() {
	app.Logger.Sync()
}
