package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/anoriar/gophkeeper/internal/server/entry/enum"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"github.com/anoriar/gophkeeper/internal/server/entry/dto/collection"
	"github.com/anoriar/gophkeeper/internal/server/entry/entity"
	"github.com/anoriar/gophkeeper/internal/server/shared/app/db"
	customCtx "github.com/anoriar/gophkeeper/internal/server/shared/context"
	errors2 "github.com/anoriar/gophkeeper/internal/server/shared/errors"
)

type EntryRepository struct {
	db *db.Database
}

func NewEntryRepository(db *db.Database) *EntryRepository {
	return &EntryRepository{db: db}
}

func (e *EntryRepository) GetEntriesByUserIDAndType(ctx context.Context, userID string, entryType enum.EntryType) (collection.EntryCollection, error) {
	var resultItems []entity.Entry
	rows, err := e.db.Conn.QueryxContext(ctx, "SELECT * FROM entries WHERE user_id = $1 AND type = $2", userID, string(entryType))
	if err != nil {
		return *collection.NewEntryCollection(nil), fmt.Errorf("GetEntriesByUserIDAndType: %w: %v", errors2.ErrInternalError, err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry entity.Entry
		err := rows.StructScan(&entry)
		if err != nil {
			return *collection.NewEntryCollection(nil), fmt.Errorf("GetEntriesByUserIDAndType: %w: %v", errors2.ErrInternalError, err)
		}
		resultItems = append(resultItems, entry)
	}

	if rows.Err() != nil {
		return *collection.NewEntryCollection(nil), fmt.Errorf("GetEntriesByUserIDAndType: %w: %v", errors2.ErrInternalError, err)
	}

	return *collection.NewEntryCollection(resultItems), nil
}

func (e *EntryRepository) AddEntries(ctx context.Context, entries []entity.Entry) error {
	txx, err := e.getTxFromContextOrBeginNew(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
	}

	stmt, err := txx.PreparexContext(ctx, "INSERT INTO entries (id, type, user_id, updated_at, data, meta, original_id) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		return fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
	}

	for _, entry := range entries {
		_, err := stmt.ExecContext(ctx, entry.Id, entry.EntryType, entry.UserId, entry.UpdatedAt, entry.Data, entry.Meta, entry.OriginalId)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
				return fmt.Errorf("%w: %v", errors2.ErrConflict, err)
			}
			return fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
		}
	}

	return nil
}

func (e *EntryRepository) UpdateEntries(ctx context.Context, entries []entity.Entry) error {
	txx, err := e.getTxFromContextOrBeginNew(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
	}

	stmt, err := txx.PreparexContext(ctx, "UPDATE entries SET type = $1, user_id = $2, updated_at = $3, data = $4, meta = $5, original_id = $6 WHERE id = $7")
	if err != nil {
		return fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
	}

	for _, entry := range entries {
		_, err := stmt.ExecContext(ctx, entry.EntryType, entry.UserId, entry.UpdatedAt, entry.Data, entry.Meta, entry.OriginalId, entry.Id)
		if err != nil {
			return err
		}
		if err != nil {
			return fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
		}
	}

	return nil
}

func (e *EntryRepository) DeleteEntries(ctx context.Context, entriesIds []string) error {
	txx, err := e.getTxFromContextOrBeginNew(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
	}

	stmt, err := txx.PreparexContext(ctx, "DELETE FROM entries WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
	}
	for _, entryId := range entriesIds {
		_, err := stmt.ExecContext(ctx, entryId)
		if err != nil {
			return err
		}
		if err != nil {
			return fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
		}
	}

	return nil
}

func (e *EntryRepository) getTxFromContextOrBeginNew(ctx context.Context) (*sqlx.Tx, error) {
	var txx *sqlx.Tx
	var err error
	txCtxParam := ctx.Value(customCtx.TransactionKey)
	if txCtxParam != nil {
		tx, ok := txCtxParam.(*db.DBTransaction)
		if ok {
			if txx, ok = tx.GetTransaction().(*sqlx.Tx); ok {
				return txx, nil
			} else {
				return nil, fmt.Errorf("%w: %v", errors2.ErrInternalError, "can not get transaction")
			}
		} else {
			return nil, fmt.Errorf("%w: %v", errors2.ErrInternalError, "can not get transaction")
		}
	} else {
		txx, err = e.db.Conn.BeginTxx(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errors2.ErrInternalError, err)
		}
	}
	return txx, nil
}
