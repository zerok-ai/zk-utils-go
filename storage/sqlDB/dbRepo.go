package sqlDB

import (
	"database/sql"
	"github.com/zerok-ai/zk-utils-go/interfaces"
)

type DatabaseRepo interface {
	//Below are the methods for transaction
	CreateTransaction() (*sql.Tx, error)
	CreateTransactionWithIsolation(isolation sql.IsolationLevel) (*sql.Tx, error)
	CommitTransaction(tx *sql.Tx) error
	RollbackTransaction(tx *sql.Tx) error
	GetAllWithTx(tx *sql.Tx, query string, param []any) (*sql.Rows, error, func())
	GetWithTx(tx *sql.Tx, query string, param []any, args []any) error

	Get(query string, param []any, args []any) error
	GetAll(query string, param []any) (*sql.Rows, error, func())
	Insert(stmt *sql.Stmt, data interfaces.DbArgs) (sql.Result, error)
	InsertWithReturnRow(stmt *sql.Stmt, param []any, args []any) error
	Update(stmt *sql.Stmt, param []any) (sql.Result, error)
	Delete(stmt *sql.Stmt, param []any) (sql.Result, error)
	Upsert(stmt *sql.Stmt, data interfaces.DbArgs) (sql.Result, error)
	BulkInsertUsingCopyIn(stmt *sql.Stmt, data []interfaces.DbArgs) error
	BulkUpsert(stmt *sql.Stmt, data []interfaces.DbArgs) ([]sql.Result, error)
	Close() error
	CreateStatement(query string) *sql.Stmt
}
