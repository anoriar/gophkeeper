package sync

import (
	"context"
	"github.com/anoriar/gophkeeper/internal/server/storage/dto/request"
	"github.com/anoriar/gophkeeper/internal/server/storage/entity"
)

type SyncServiceInterface interface {
	Sync(ctx context.Context, request request.SyncRequest) ([]entity.Entry, error)
}
