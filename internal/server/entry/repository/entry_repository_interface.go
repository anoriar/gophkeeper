package repository

import (
	"context"

	"github.com/anoriar/gophkeeper/internal/server/entry/enum"

	"github.com/anoriar/gophkeeper/internal/server/entry/dto/collection"
	"github.com/anoriar/gophkeeper/internal/server/entry/entity"
)

//go:generate mockgen -source=entry_repository_interface.go -destination=entry_repository_mock/entry_repository.go -package=entry_repository_mock
type EntryRepositoryInterface interface {
	GetEntriesByUserIDAndType(ctx context.Context, userID string, entryType enum.EntryType) (collection.EntryCollection, error)
	AddEntries(ctx context.Context, entries []entity.Entry) error
	UpdateEntries(ctx context.Context, entries []entity.Entry) error
	DeleteEntries(ctx context.Context, entriesIds []string) error
}
