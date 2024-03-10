package entry_ext

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/repository/entry_ext"
	entryErr "github.com/anoriar/gophkeeper/internal/client/entry/errors"
	sharedErr "github.com/anoriar/gophkeeper/internal/client/shared/errors"
)

type EntryExtRepository struct {
	client *resty.Client
}

func NewEntryExtRepository(client *resty.Client) *EntryExtRepository {
	return &EntryExtRepository{client: client}
}

func (e *EntryExtRepository) Sync(ctx context.Context, token string, entries []entry_ext.SyncRequestItem) (entry_ext.SyncResponse, error) {
	body, err := json.Marshal(map[string]interface{}{"items": entries})
	if err != nil {
		return entry_ext.SyncResponse{}, fmt.Errorf("%w: %v", sharedErr.ErrInternalError, err)
	}

	resp, err := e.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", token).
		SetBody(body).
		Post("/api/entries/sync")

	if err != nil {
		return entry_ext.SyncResponse{}, fmt.Errorf("%w: %v", sharedErr.ErrDependencyFailure, err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		var result entry_ext.SyncResponse
		err = json.Unmarshal(resp.Body(), &result)
		if err != nil {
			return entry_ext.SyncResponse{}, fmt.Errorf("%w: %v", sharedErr.ErrInternalError, err)
		}
		return result, nil
	case http.StatusConflict:
		return entry_ext.SyncResponse{}, fmt.Errorf("%w: %v", entryErr.ErrSyncConflict, resp.Body())
	default:
		return entry_ext.SyncResponse{}, fmt.Errorf("%w: %v", sharedErr.ErrDependencyFailure, resp.Body())
	}
}
