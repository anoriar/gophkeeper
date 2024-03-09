package sync

import (
	"encoding/json"
	"errors"
	sharedErrors "github.com/anoriar/gophkeeper/internal/server/shared/errors"
	"io"
	"net/http"

	entryErrors "github.com/anoriar/gophkeeper/internal/server/entry/errors"

	"go.uber.org/zap"

	sync2 "github.com/anoriar/gophkeeper/internal/server/entry/dto/request/sync"
	"github.com/anoriar/gophkeeper/internal/server/entry/services/sync"
	customCtx "github.com/anoriar/gophkeeper/internal/server/shared/context"
)

type SyncHandler struct {
	syncService sync.SyncServiceInterface
	logger      *zap.Logger
}

func NewSyncHandler(syncService sync.SyncServiceInterface, logger *zap.Logger) *SyncHandler {
	return &SyncHandler{syncService: syncService, logger: logger}
}

func (sh *SyncHandler) Sync(w http.ResponseWriter, req *http.Request) {
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		sh.logger.Error("internal server error", zap.String("error", err.Error()))
		return
	}

	var syncRequest sync2.SyncRequest
	err = json.Unmarshal(requestBody, &syncRequest)

	if err != nil {
		if _, ok := err.(*json.SyntaxError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if errors.Is(err, entryErrors.ErrSyncRequestNotValid) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			sh.logger.Error("unmarshal error", zap.String("error", err.Error()))
		}
		return
	}
	userID := ""
	userIDCtxParam := req.Context().Value(customCtx.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}

	if userID == "" {
		http.Error(w, "user unauthorized", http.StatusUnauthorized)
		return
	}

	syncRequest.UserID = userID

	response, err := sh.syncService.Sync(req.Context(), syncRequest)
	if err != nil {
		switch {
		case errors.Is(err, entryErrors.ErrSyncRequestNotValid):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, sharedErrors.ErrConflict):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			sh.logger.Error("internal server error", zap.String("error", err.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResult, err := json.Marshal(response)
	if err != nil {
		sh.logger.Error("marshal error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResult)
	if err != nil {
		sh.logger.Error("write command_response error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
