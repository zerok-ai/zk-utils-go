package sqlDB

import (
	"database/sql"
	"github.com/zerok-ai/zk-utils-go/interfaces"
)

type DatabaseRepo interface {
	Get(query string, param []any, args []any) error
	GetAll(query string, param []any) (*sql.Rows, error, func())
	Insert(stmt *sql.Stmt, data interfaces.DbArgs) (sql.Result, error)
	InsertWithReturnRow(stmt *sql.Stmt, param []any) (*sql.Row, error)
	Update(stmt *sql.Stmt, param []any) (sql.Result, error)
	Delete(stmt *sql.Stmt, param []any) (sql.Result, error)
	Upsert(stmt *sql.Stmt, data interfaces.DbArgs) (sql.Result, error)
	BulkInsertUsingCopyIn(stmt *sql.Stmt, data []interfaces.DbArgs) error
	BulkUpsert(stmt *sql.Stmt, data []interfaces.DbArgs) ([]sql.Result, error)
	CreateTransaction() (*sql.Tx, error)
	Close() error
}
