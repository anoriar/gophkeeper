package db

//go:generate mockgen -source=db_transaction_interface.go -destination=mock/db_transaction.go -package=mock
type DBTransactionInterface interface {
	Commit() error
	Rollback() error
	GetTransaction() interface{}
}
