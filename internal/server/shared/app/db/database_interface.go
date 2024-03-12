package db

import "context"

//go:generate mockgen -source=database_interface.go -destination=mock/database.go -package=mock
type DatabaseInterface interface {
	Ping() error
	BeginTransaction(ctx context.Context) (DBTransactionInterface, error)
	Close()
}
