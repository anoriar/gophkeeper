package entry_ext

import (
	"context"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/repository/entry_ext"
)

//go:generate mockgen -source=entry_ext_repository_interface.go -destination=mock_entry_ext_repository/mock_entry_ext_repository.go -package=mock_entry_ext_repository
type EntryExtRepositoryInterface interface {
	Sync(ctx context.Context, token string, request entry_ext.SyncRequest) (entry_ext.SyncResponse, error)
}
