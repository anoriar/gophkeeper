package auth

import (
	"context"
	"errors"
	"net/http"

	sharedErrors "github.com/anoriar/gophkeeper/internal/server/shared/errors"

	context2 "github.com/anoriar/gophkeeper/internal/server/shared/context"
	auth2 "github.com/anoriar/gophkeeper/internal/server/user/services/auth"
)

type AuthMiddleware struct {
	authService auth2.AuthServiceInterface
}

func NewAuthMiddleware(authService auth2.AuthServiceInterface) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (middleware *AuthMiddleware) Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		token := request.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "user unauthorized", http.StatusUnauthorized)
			return
		}
		claims, err := middleware.authService.ValidateToken(token)
		if err != nil {
			if errors.Is(err, sharedErrors.ErrUserUnauthorized) {
				http.Error(w, "user unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(request.Context(), context2.UserIDContextKey, claims.UserID)

		h.ServeHTTP(w, request.WithContext(ctx))
	})
}
