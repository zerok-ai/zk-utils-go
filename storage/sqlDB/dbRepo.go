package sqlDB

import (
	"database/sql"
	"github.com/zerok-ai/zk-utils-go/interfaces"
)

type DatabaseRepo interface {
	Get(query string, param []any, args []any) error
	GetAll(query string, param []any) (*sql.Rows, error, func())
	//Insert(query string, data interfaces.DbArgs) (sql.Result, error)
	BulkInsert(stmt *sql.Stmt, data []interfaces.DbArgs) error
	Insert(stmt *sql.Stmt, data interfaces.DbArgs) (sql.Result, error)
	Update(stmt *sql.Stmt, param []any) (int, error)
	BulkUpsert(stmt *sql.Stmt, data []interfaces.DbArgs) error
	Delete(stmt *sql.Stmt, param []any) (int, error)
	CreateTransaction() (*sql.Tx, error)
}
