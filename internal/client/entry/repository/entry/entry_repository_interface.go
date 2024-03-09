package entry

import (
	"context"

	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
)

type EntryRepositoryInterface interface {
	Add(ctx context.Context, entry entity.Entry) error
	Edit(ctx context.Context, entry entity.Entry) error
	GetById(ctx context.Context, id string) (entity.Entry, error)
	GetList(ctx context.Context) ([]entity.Entry, error)
	Rewrite(ctx context.Context, entries []entity.Entry) error
}
