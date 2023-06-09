package sqlDB

import (
	"database/sql"
	"github.com/zerok-ai/zk-utils-go/interfaces"
)

type DatabaseRepo interface {
	GetDBInstance() (*sql.DB, error)
	Delete(tx *sql.Tx, query string, param []any) (int, error)
	Update(tx *sql.Tx, stmt string, param []any) (int, error)
	Get(db *sql.DB, query string, param []any, args []any) error
	GetAll(db *sql.DB, query string, param []any) (*sql.Rows, error, func())
	Insert(db *sql.DB, query string, param []any) error
	BulkInsert(tx *sql.Tx, tableName string, columns []string, data []interfaces.DbArgs) error
	InsertInTransaction(tx *sql.Tx, stmt string, params []any) error
	BulkUpsert(tx *sql.Tx, stmt string, data [][]interface{}) error
}
