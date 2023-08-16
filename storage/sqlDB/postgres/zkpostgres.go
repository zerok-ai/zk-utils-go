package zkpostgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/zerok-ai/zk-utils-go/interfaces"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB"
	pgConfig "github.com/zerok-ai/zk-utils-go/storage/sqlDB/postgres/config"
	"log"
	"time"
)

// A new type is created for postgres, this will implement the interface DatabaseRepo.
type zkPostgresRepo struct {
	Db *sql.DB
}

func (databaseRepo zkPostgresRepo) InsertWithReturnRow(stmt *sql.Stmt, param []any, args []any) error {
	if stmt == nil {
		err := errors.New("statement cannot be empty")
		zkLogger.Error(LogTag, err)
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(param...)
	return row.Scan(args...)
}

var LogTag = "zkpostgres_db_repo"
var config pgConfig.PostgresConfig
var dbInstance *sql.DB

// NewZkPostgresRepo This method returns an instance of DatabaseRepo type where the underlying type is zkPostgresRepo.
// This instance will be used to call the Postgres implementation of DatabaseRepo methods.
// If in PostgresConfig no value is passed for MaxConnections, MaxIdleConnections or ConnectionMaxLifetimeInMinutes,
// the default value of them is 10, 5 and 30 minutes respectively
func NewZkPostgresRepo(c pgConfig.PostgresConfig) (sqlDB.DatabaseRepo, error) {
	config = c
	instance, err := getDBInstance()
	zkRepo := zkPostgresRepo{
		Db: instance,
	}
	return zkRepo, err
}

// a non exported function to create the connection poll for Postgres
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

// if no value is passed for MaxConnections, MaxIdleConnections or ConnectionMaxLifetimeInMinutes,
// the default value of them is 10, 5 and 30 minutes respectively
func getDBInstance() (*sql.DB, error) {
	var err error

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

	return dbInstance, err
}

// Get This method takes a db handle, sql query, slice of params where each value in the param provides values to
// placeholders in the query. For example in the query:
// SELECT name, age FROM STUDENTS WHERE id=$1, the param slice will have just 1 value for id.
// The fourth arguments is a slice of arguments where each value in args corresponds to the columns in select query, Ex:
// in query: SELECT name, age FROM STUDENTS WHERE id=$1, args will have 2 values
// args[0] = &student.name and args[1] = &student.age, so that the values are persisted in the student object after the
// function call ends as we are not returning anything other than error
func (databaseRepo zkPostgresRepo) Get(query string, param []any, args []any) error {
	if query == "" {
		err := errors.New("query cannot be empty")
		zkLogger.Error(LogTag, err)
		return err
	}

	db := databaseRepo.Db
	defer db.Close()
	row := db.QueryRow(query, param...)
	return row.Scan(args...)
}

// GetAll This method takes a db handle, sql query, slice of params where each value in the param provides values to
// placeholders in the query. For example in the query:
// SELECT name, age FROM STUDENTS WHERE id=$1, the param slice will have just 1 value for id.
// we cannot pass []args as we did in Get since the number of rows returned from select query is unknown, so we cannot
// initialize a slice so that we can set values into its variables.
// That's why we are returning Rows and handling the rows inside the code.
func (databaseRepo zkPostgresRepo) GetAll(query string, param []any) (*sql.Rows, error, func()) {
	if query == "" {
		err := errors.New("query cannot be empty")
		zkLogger.Error(LogTag, err)
		return nil, err, func() {}
	}
	db := databaseRepo.Db
	rows, err := db.Query(query, param...)
	closeRow := func() {
		defer rows.Close()
	}
	return rows, err, closeRow
}

// Insert This method takes a db handle, sql query and DbArgs
// DbArgs is an interface which defines a method GetAllColumns which returns a slice of struct fields corresponding to columns
// in db table, The values from this slice is then inserted into the table.
func (databaseRepo zkPostgresRepo) Insert(stmt *sql.Stmt, data interfaces.DbArgs) (sql.Result, error) {
	return databaseRepo.modifyTable(stmt, data.GetAllColumns())
}

// Update This method takes a transaction, sql query, slice of params where each value in the param provides values to
// placeholders in the query. For more details of placeholder refer Get or GetAll
// it returns number of rows modified and error
func (databaseRepo zkPostgresRepo) Update(stmt *sql.Stmt, param []any) (sql.Result, error) {
	return databaseRepo.modifyTable(stmt, param)
}

// Delete This method takes a transaction, sql query, slice of params where each value in the param provides values to
// placeholders in the query. For more details of placeholder refer Get or GetAll
// it returns number of rows modified and error
func (databaseRepo zkPostgresRepo) Delete(stmt *sql.Stmt, param []any) (sql.Result, error) {
	return databaseRepo.modifyTable(stmt, param)
}

// Upsert This method takes a transaction, sql query and a slice DbArgs
// For more details on DbArgs read Insert
// Here we are using Upsert command which gives additional control when inserts operation fail due to conflicts.
// Example Upsert Query: INSERT INTO TableA (colA, colB) VALUES ($1, $2) ON CONFLICT (colA) DO NOTHING
// The above query tries to insert some value to TableA but on conflict, it does nothing, you can also specify what to do
// incase of a conflict.
func (databaseRepo zkPostgresRepo) Upsert(stmt *sql.Stmt, data interfaces.DbArgs) (sql.Result, error) {
	return databaseRepo.modifyTable(stmt, data.GetAllColumns())
}

// BulkInsertUsingCopyIn This method takes a transaction, tableName, list of columns and corresponding column values in []DbArgs
// For more details on DbArgs read Insert
// Here we are using copyIn which is a very fast operation for bulk Insert
func (databaseRepo zkPostgresRepo) BulkInsertUsingCopyIn(stmt *sql.Stmt, data []interfaces.DbArgs) error {
	if stmt == nil {
		err := errors.New("statement cannot be empty")
		zkLogger.Error(LogTag, err)
		return err
	}

	defer stmt.Close()
	var results []sql.Result

	for _, d := range data {
		c := d.GetAllColumns()
		r, err := stmt.Exec(c...)
		results = append(results, r)
		if err != nil {
			zkLogger.Error(LogTag, "couldn't prepare COPY statement: %v", err)
			return err
		}
	}

	finalResult, err := stmt.Exec()
	results = append(results, finalResult)

	if err != nil {
		zkLogger.Error(LogTag, "failed to perform bulk insert using copyIn: ", err)
		return err
	}

	return nil
}

// BulkUpsert Uses Upsert command to insert or update data. One can also pass ON CONFLICT DO NOTHING in the query and
// no update will happen if there is a conflict
func (databaseRepo zkPostgresRepo) BulkUpsert(stmt *sql.Stmt, data []interfaces.DbArgs) ([]sql.Result, error) {
	if stmt == nil {
		err := errors.New("statement cannot be empty")
		zkLogger.Error(LogTag, err)
		return nil, err
	}

	defer stmt.Close()
	var results []sql.Result

	for _, d := range data {
		c := d.GetAllColumns()
		r, err := stmt.Exec(c...)
		results = append(results, r)
		if err != nil {
			zkLogger.Error(LogTag, "failed to perform bulk upsert: ", err)
			return results, err
		}
	}

	return results, nil
}

// CreateTransaction Returns a new transaction
func (databaseRepo zkPostgresRepo) CreateTransaction() (*sql.Tx, error) {
	ctx := context.Background()
	tx, err := databaseRepo.Db.BeginTx(ctx, nil)
	if err != nil {
		zkLogger.Debug(LogTag, "unable to create txn, "+err.Error())
		return nil, err
	}
	return tx, nil
}

// internal method used by update and delete
func (databaseRepo zkPostgresRepo) modifyTable(stmt *sql.Stmt, param []any) (sql.Result, error) {
	if stmt == nil {
		err := errors.New("statement cannot be empty")
		zkLogger.Error(LogTag, err)
		return nil, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(param...)
	if err == nil {
		return res, nil
	}

	zkLogger.Debug(LogTag, err.Error())
	return nil, err
}

func (databaseRepo zkPostgresRepo) Close() error {
	return databaseRepo.Db.Close()
}
