package sqlDB

import (
	"database/sql"
	"github.com/zerok-ai/zk-utils-go/interfaces"
)

type DatabaseRepo interface {
	Get(query string, param []any, args []any) error
	GetAll(query string, param []any) (*sql.Rows, error, func())
	Insert(query string, data interfaces.DbArgs) (sql.Result, error)
	BulkInsert(tableName string, columns []string, data []interfaces.DbArgs) error
	InsertInTransaction(stmt string, data interfaces.DbArgs) (sql.Result, error)
	Update(stmt string, param []any) (int, error)
	BulkUpsert(stmt string, data []interfaces.DbArgs) error
	Delete(query string, param []any) (int, error)
}
