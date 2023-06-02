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
)

type zkPostgresRepo struct {
}

func NewZkPostgresRepo() sqlDB.DatabaseRepo {
	return &zkPostgresRepo{}
}

var config pgConfig.PostgresConfig
var LogTag = "zkpostgres_db_repo"

func Init(c pgConfig.PostgresConfig) {
	config = c
}

func (zkPostgresService zkPostgresRepo) CreateConnection() *sql.DB {
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

func (zkPostgresService zkPostgresRepo) Get(db *sql.DB, query string, param []any, args []any) error {
	defer db.Close()
	row := db.QueryRow(query, param...)
	return row.Scan(args...)
}

func (zkPostgresService zkPostgresRepo) GetAll(db *sql.DB, query string, param []any) (*sql.Rows, error, func()) {
	rows, err := db.Query(query, param...)
	closeRow := func() {
		defer rows.Close()
	}
	return rows, err, closeRow
}

func (zkPostgresService zkPostgresRepo) Delete(tx *sql.Tx, query string, param []any) (int, error) {
	return zkPostgresService.modifyTable(tx, query, param)
}

func (zkPostgresService zkPostgresRepo) Update(tx *sql.Tx, stmt string, param []any) (int, error) {
	return zkPostgresService.modifyTable(tx, stmt, param)
}

func (zkPostgresService zkPostgresRepo) modifyTable(tx *sql.Tx, stmt string, param []any) (int, error) {
	preparedStmt, err := tx.Prepare(stmt)
	if err != nil {
		zkLogger.Error(LogTag, "Error preparing delete/update statement:", err)
		return 0, err
	}

	defer preparedStmt.Close()
	res, err := preparedStmt.Exec(param...)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			return int(count), nil
		} else {
			zkLogger.Error(LogTag, "Error executing update/delete:", err)
			return 0, err
		}
	}

	zkLogger.Debug(LogTag, err.Error())
	return 0, err
}

func (zkPostgresService zkPostgresRepo) Insert(db *sql.DB, query string, param []any) error {
	preparedStmt, err := db.Prepare(query)
	if err != nil {
		zkLogger.Error(LogTag, "Error preparing insert statement:", err)
		return err
	}
	defer preparedStmt.Close()

	_, err = preparedStmt.Exec(param...)
	if err != nil {
		zkLogger.Error(LogTag, "Error executing insert statement:", err)
		return err
	}

	return nil
}

func (zkPostgresService zkPostgresRepo) BulkInsert(tx *sql.Tx, tableName string, columns []string, data []interfaces.DbArgs) error {
	stmt, err := tx.Prepare(pq.CopyIn(tableName, columns...))
	if err != nil {
		return err
	}

	for _, d := range data {
		values := d.GetArgs()
		_, err := stmt.Exec(values...)
		if err != nil {
			zkLogger.Error(LogTag, "couldn't prepare COPY statement: %v", err)
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		zkLogger.Error(LogTag, "Error executing copy statement:", err)
		return err
	}

	return stmt.Close()
}

func (zkPostgresService zkPostgresRepo) InsertInTransaction(tx *sql.Tx, stmt string, params []any) error {
	preparedStmt, err := tx.Prepare(stmt)
	if err != nil {
		zkLogger.Error(LogTag, "Error preparing insert statement:", err)
		return err
	}
	defer preparedStmt.Close()

	_, err = preparedStmt.Exec(params...)
	if err != nil {
		zkLogger.Error(LogTag, "Error executing insert:", err)
	}

	return nil
}
