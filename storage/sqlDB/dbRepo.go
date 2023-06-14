package sqlDB

import (
	"database/sql"
	"github.com/zerok-ai/zk-utils-go/interfaces"
)

type DatabaseRepo interface {
	Get(query string, param []any, args []any) error
	GetAll(query string, param []any) (*sql.Rows, error, func())
	//Insert(query string, data interfaces.DbArgs) (sql.Result, error)
	Insert(stmt *sql.Stmt, data interfaces.DbArgs) (sql.Result, error)
	Update(stmt *sql.Stmt, param []any) (int, error)
	Delete(stmt *sql.Stmt, param []any) (int, error)
	Upsert(stmt *sql.Stmt, data interfaces.DbArgs) error
	BulkInsertUsingCopyIn(stmt *sql.Stmt, data []interfaces.DbArgs) error
	BulkUpsert(stmt *sql.Stmt, data []interfaces.DbArgs) error
	CreateTransaction() (*sql.Tx, error)
}
