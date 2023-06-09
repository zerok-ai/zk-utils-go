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
	"sync"
	"time"
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

var (
	dbInstance *sql.DB
	once       sync.Once
)

func createConnectionPool(connectionString string, maxConnections int, maxIdleConnections int, connectionMaxLifetime time.Duration) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		zkLogger.Error("failed to open database connection database:", err)
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	db.SetMaxOpenConns(maxConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxLifetime(connectionMaxLifetime)

	err = db.Ping()
	if err != nil {
		zkLogger.Error("failed to ping database:", err)
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

func (zkPostgresService zkPostgresRepo) GetDBInstance() (*sql.DB, error) {
	var err error

	once.Do(func() {
		connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.User, config.Password, config.Dbname)

		var maxConnections, maxIdleConnections int
		var connectionMaxLifetime time.Duration

		if config.MaxConnections == 0 {
			maxConnections = 10
		}
		if config.MaxIdleConnections == 0 {
			maxIdleConnections = 5
		}
		if config.ConnectionMaxLifetimeInMinutes == 0 {
			connectionMaxLifetime = time.Minute * 30
		}

		// Create the connection pool
		dbInstance, err = createConnectionPool(connectionString, maxConnections, maxIdleConnections, connectionMaxLifetime)
		if err != nil {
			log.Fatalf("failed to create connection pool: %v", err)
		}
	})

	return dbInstance, err
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
		zkLogger.Error(LogTag, "Error preparing insert statement:", err)
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
		return err
	}

	return nil
}

func (zkPostgresService zkPostgresRepo) BulkUpsert(tx *sql.Tx, stmt string, data []interfaces.DbArgs) error {
	preparedStmt, err := tx.Prepare(stmt)
	if err != nil {
		zkLogger.Error(LogTag, "Error preparing insert statement:", err)
		return err
	}

	defer preparedStmt.Close()

	for _, row := range data {
		_, err = preparedStmt.Exec(row.GetArgs())
		if err != nil {
			_ = tx.Rollback()
			zkLogger.Error(LogTag, "failed to perform bulk upsert: ", err)
			return err
		}
	}

	return nil
}
