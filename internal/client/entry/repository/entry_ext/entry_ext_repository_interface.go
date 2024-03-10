package entry_ext

import (
	"context"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/repository/entry_ext"
)

type EntryExtRepositoryInterface interface {
	Sync(ctx context.Context, token string, entries []entry_ext.SyncRequestItem) (entry_ext.SyncResponse, error)
}
