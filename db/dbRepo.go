package db

import "database/sql"

type DatabaseRepo[T any] interface {
	CreateConnection() *sql.DB
	Delete(stmt string, param []any, tx *sql.Tx, rollback bool) (int, error)
	Update(stmt string, param []any, tx *sql.Tx, rollback bool) (int, error)
	Get(query string, param []any, args []any) error
	Insert(stmt string, param []any) error
}
