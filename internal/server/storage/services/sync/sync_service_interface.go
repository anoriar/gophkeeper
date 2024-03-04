package sync

import (
	"context"

	"github.com/anoriar/gophkeeper/internal/server/storage/dto/request/sync"
	syncResponsePkg "github.com/anoriar/gophkeeper/internal/server/storage/dto/response/sync"
)

type SyncServiceInterface interface {
	Sync(ctx context.Context, request sync.SyncRequest) (syncResponsePkg.SyncResponse, error)
}
