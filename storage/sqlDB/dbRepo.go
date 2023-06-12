package sqlDB

import (
	"database/sql"
	"github.com/zerok-ai/zk-utils-go/interfaces"
)

type DatabaseRepo interface {
	GetDBInstance() (*sql.DB, error)
	Get(db *sql.DB, query string, param []any, args []any) error
	GetAll(db *sql.DB, query string, param []any) (*sql.Rows, error, func())
	Insert(db *sql.DB, query string, data interfaces.DbArgs) (sql.Result, error)
	BulkInsert(tx *sql.Tx, tableName string, columns []string, data []interfaces.DbArgs) error
	InsertInTransaction(tx *sql.Tx, stmt string, data interfaces.DbArgs) (sql.Result, error)
	Update(tx *sql.Tx, stmt string, param []any) (int, error)
	BulkUpsert(tx *sql.Tx, stmt string, data []interfaces.DbArgs) error
	Delete(tx *sql.Tx, query string, param []any) (int, error)
}
