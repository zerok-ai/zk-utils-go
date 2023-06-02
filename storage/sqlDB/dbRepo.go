package sqlDB

import (
	"database/sql"
	"github.com/zerok-ai/zk-utils-go/interfaces"
)

type DatabaseRepo[T any] interface {
	CreateConnection() *sql.DB
	Delete(tx *sql.Tx, query string, param []any, rollback bool) (int, error)
	Update(tx *sql.Tx, stmt string, param []any, rollback bool) (int, error)
	Get(db *sql.DB, query string, param []any, args []any) error
	GetAll(db *sql.DB, query string, param []any) (*sql.Rows, error, func())
	Insert(db *sql.DB, query string, param []any) error
	BulkInsert(tx *sql.Tx, tableName string, columns []string, data []interfaces.PostgresRuleIterator) error
}
