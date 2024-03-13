package entry

import (
	"context"

	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
)

//go:generate mockgen -source=entry_repository_interface.go -destination=mock_entry_repository/mock_entry_repository.go -package=mock_entry_repository
type EntryRepositoryInterface interface {
	Add(ctx context.Context, entry entity.Entry) error
	Edit(ctx context.Context, entry entity.Entry) error
	GetById(ctx context.Context, id string) (entity.Entry, error)
	GetList(ctx context.Context) ([]entity.Entry, error)
	Rewrite(ctx context.Context, entries []entity.Entry) error
}
