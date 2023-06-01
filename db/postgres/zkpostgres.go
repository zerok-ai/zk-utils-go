package zkpostgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/zerok-ai/zk-utils-go/db"
	pgConfig "github.com/zerok-ai/zk-utils-go/db/postgres/config"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"log"
)

type zkPostgresRepo[T any] struct {
}

func NewZkPostgresRepo[T any]() db.DatabaseRepo[T] {
	return &zkPostgresRepo[T]{}
}

var config pgConfig.PostgresConfig
var LogTag = "zkpostgres_db_repo"

func Init(c pgConfig.PostgresConfig) {
	config = c
}

func (zkPostgresService zkPostgresRepo[T]) CreateConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)
	zkLogger.Debug(LogTag, "psqlInfo==", psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func (zkPostgresService zkPostgresRepo[T]) Get(query string, param []any, args []any) error {
	db := zkPostgresService.CreateConnection()
	defer db.Close()
	row := db.QueryRow(query, param...)
	return row.Scan(args...)
}

//func (zkPostgresService zkPostgresRepo[T]) rowProcessor(row *sql.Row, args []any) error {
//	return row.Scan(args...)
//
//	zkErr := zkcommon.CheckSqlError(err, LOG_TAG)
//	if zkErr != nil {
//		return zkErr
//	}
//	return nil
//}

func (zkPostgresService zkPostgresRepo[T]) Delete(stmt string, param []any, tx *sql.Tx, rollback bool) (int, error) {
	return zkPostgresService.modifyTable(stmt, param, tx, rollback)
}

func (zkPostgresService zkPostgresRepo[T]) Update(stmt string, param []any, tx *sql.Tx, rollback bool) (int, error) {
	return zkPostgresService.modifyTable(stmt, param, tx, rollback)
}

func (zkPostgresService zkPostgresRepo[T]) modifyTable(stmt string, param []any, tx *sql.Tx, rollback bool) (int, error) {
	res, err := tx.Exec(stmt, param...)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			return int(count), nil
		} else {
			return 0, err
		}
	}

	zkLogger.Debug(LogTag, err.Error())
	return 0, err
}

func (zkPostgresService zkPostgresRepo[T]) Insert(stmt string, param []any) error {
	db := zkPostgresService.CreateConnection()
	preparedStmt, err := db.Prepare(stmt)
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer preparedStmt.Close()

	_, err = preparedStmt.Exec(param...)
	if err != nil {
		log.Println("Error executing statement:", err)
		return err
	}

	return nil
}
