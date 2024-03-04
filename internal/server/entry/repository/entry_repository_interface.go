package repository

import (
	"context"

	"github.com/anoriar/gophkeeper/internal/server/entry/dto/collection"
	"github.com/anoriar/gophkeeper/internal/server/entry/entity"
)

type EntryRepositoryInterface interface {
	GetEntriesByUserID(ctx context.Context, userID string) (collection.EntryCollection, error)
	AddEntries(ctx context.Context, entries []entity.Entry) error
	UpdateEntries(ctx context.Context, entries []entity.Entry) error
	DeleteEntries(ctx context.Context, entriesIds []string) error
}
