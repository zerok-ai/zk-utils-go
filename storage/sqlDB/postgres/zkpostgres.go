package zkpostgres

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/zerok-ai/zk-utils-go/interfaces"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB"
	pgConfig "github.com/zerok-ai/zk-utils-go/storage/sqlDB/postgres/config"
	"log"
)

type zkPostgresRepo[T any] struct {
}

func NewZkPostgresRepo[T any]() sqlDB.DatabaseRepo[T] {
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

func (zkPostgresService zkPostgresRepo[T]) Get(db *sql.DB, query string, param []any, args []any) error {
	defer db.Close()
	row := db.QueryRow(query, param...)
	return row.Scan(args...)
}

func (zkPostgresService zkPostgresRepo[T]) GetAll(db *sql.DB, query string, param []any) (*sql.Rows, error, func()) {

	rows, err := db.Query(query, param...)
	f := func() {
		defer rows.Close()
	}
	return rows, err, f
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

func (zkPostgresService zkPostgresRepo[T]) Delete(tx *sql.Tx, query string, param []any, rollback bool) (int, error) {
	return zkPostgresService.modifyTable(tx, query, param, rollback)
}

func (zkPostgresService zkPostgresRepo[T]) Update(tx *sql.Tx, stmt string, param []any, rollback bool) (int, error) {
	return zkPostgresService.modifyTable(tx, stmt, param, rollback)
}

func (zkPostgresService zkPostgresRepo[T]) modifyTable(tx *sql.Tx, stmt string, param []any, rollback bool) (int, error) {
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

func (zkPostgresService zkPostgresRepo[T]) Insert(db *sql.DB, query string, param []any) error {
	preparedStmt, err := db.Prepare(query)
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

func (zkPostgresService zkPostgresRepo[T]) BulkInsert(tx *sql.Tx, tableName string, columns []string, data []interfaces.PostgresRuleIterator) error {
	stmt, err := tx.Prepare(pq.CopyIn(tableName, columns...))
	if err != nil {
		return err
	}
	for _, d := range data {
		values := d.Explode()
		_, err := stmt.Exec(values...)
		if err != nil {
			zkLogger.Debug("couldn't prepare COPY statement: %v", err)
			return err
		}
		if err != nil {
			return err
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return stmt.Close()
}
