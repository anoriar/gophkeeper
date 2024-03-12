package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	Conn *sqlx.DB
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{Conn: db}
}

func (db *Database) Ping() error {
	err := db.Conn.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) BeginTransaction(ctx context.Context) (DBTransactionInterface, error) {
	txx, err := db.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	transaction := NewDBTransaction(txx)
	return transaction, nil
}

func (db *Database) Close() {
	db.Conn.Close()
}
