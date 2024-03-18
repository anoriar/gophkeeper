package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/anoriar/gophkeeper/internal/server/entry/handlers/sync"

	"github.com/anoriar/gophkeeper/internal/server/shared/app"
	"github.com/anoriar/gophkeeper/internal/server/shared/middleware/auth"
	"github.com/anoriar/gophkeeper/internal/server/shared/middleware/compress"
	"github.com/anoriar/gophkeeper/internal/server/shared/middleware/logger"
	"github.com/anoriar/gophkeeper/internal/server/user/handlers/login"
	"github.com/anoriar/gophkeeper/internal/server/user/handlers/register"
)

// Router missing godoc.
type Router struct {
	registerHandler    *register.RegisterHandler
	loginHandler       *login.LoginHandler
	syncHandler        *sync.SyncHandler
	loggerMiddleware   *logger.LoggerMiddleware
	authMiddleware     *auth.AuthMiddleware
	compressMiddleware *compress.CompressMiddleware
}

// NewRouter missing godoc.
func NewRouter(app *app.App) *Router {
	return &Router{
		registerHandler:    register.NewRegisterHandler(app.AuthService),
		loginHandler:       login.NewLoginHandler(app.AuthService),
		syncHandler:        sync.NewSyncHandler(app.SyncService, app.Logger),
		loggerMiddleware:   logger.NewLoggerMiddleware(app.Logger),
		authMiddleware:     auth.NewAuthMiddleware(app.AuthService),
		compressMiddleware: compress.NewCompressMiddleware(),
	}
}

// Route missing godoc.
func (r *Router) Route() chi.Router {
	router := chi.NewRouter()

	router.Use(r.loggerMiddleware.Log)
	router.Use(r.compressMiddleware.Compress)

	router.Post("/api/user/register", r.registerHandler.Register)
	router.Post("/api/user/login", r.loginHandler.Login)
	router.With(r.authMiddleware.Auth).Post("/api/entries/sync", r.syncHandler.Sync)

	return router
}
