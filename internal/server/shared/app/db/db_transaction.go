package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DBTransaction struct {
	tx *sqlx.Tx
}

func NewDBTransaction(tx *sqlx.Tx) *DBTransaction {
	return &DBTransaction{tx: tx}
}

func (t *DBTransaction) Commit() error {
	if t.tx == nil {
		return fmt.Errorf("no transaction in progress")
	}

	if err := t.tx.Commit(); err != nil {
		return err
	}

	t.tx = nil
	return nil
}

func (t *DBTransaction) Rollback() error {
	if t.tx == nil {
		return fmt.Errorf("no transaction in progress")
	}

	if err := t.tx.Rollback(); err != nil {
		return err
	}

	t.tx = nil
	return nil
}

func (t *DBTransaction) GetTransaction() interface{} {
	return t.tx
}
