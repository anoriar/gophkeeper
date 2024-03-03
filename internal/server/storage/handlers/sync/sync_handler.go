package sync

import (
	"encoding/json"
	customCtx "github.com/anoriar/gophkeeper/internal/server/shared/context"
	"github.com/anoriar/gophkeeper/internal/server/storage/dto/request"
	"github.com/anoriar/gophkeeper/internal/server/storage/services/sync"
	"go.uber.org/zap"
	"io"
	"net/http"
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

	var syncRequest request.SyncRequest
	err = json.Unmarshal(requestBody, &syncRequest)

	if err != nil {
		if _, ok := err.(*json.SyntaxError); ok {
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
		http.Error(w, "internal server error", http.StatusInternalServerError)
		sh.logger.Error("internal server error", zap.String("error", err.Error()))
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResult, err := json.Marshal(response)
	if err != nil {
		sh.logger.Error("marshal error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write(jsonResult)
	if err != nil {
		sh.logger.Error("write response error", zap.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
