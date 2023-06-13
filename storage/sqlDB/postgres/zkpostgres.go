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
	"time"
)

// A new type is created for postgres, this will implement the interface DatabaseRepo.
type zkPostgresRepo struct {
	Db *sql.DB
}

// NewZkPostgresRepo This method returns an instance of DatabaseRepo type where the underlying type is zkPostgresRepo.
// This instance will be used to call the Postgres implementation of DatabaseRepo methods.
func NewZkPostgresRepo(c pgConfig.PostgresConfig) (sqlDB.DatabaseRepo, error) {
	config = c
	instance, err := getDBInstance()
	if err != nil {
		return zkPostgresRepo{}, nil
	}
	zkRepo := zkPostgresRepo{
		Db: instance,
	}
	return zkRepo, nil
}

var config pgConfig.PostgresConfig
var LogTag = "zkpostgres_db_repo"

var (
	dbInstance *sql.DB
)

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
func (databaseRepo zkPostgresRepo) Get(db *sql.DB, query string, param []any, args []any) error {
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
func (databaseRepo zkPostgresRepo) GetAll(db *sql.DB, query string, param []any) (*sql.Rows, error, func()) {
	rows, err := db.Query(query, param...)
	closeRow := func() {
		defer rows.Close()
	}
	return rows, err, closeRow
}

// Delete This method takes a transaction, sql query, slice of params where each value in the param provides values to
// placeholders in the query. For more details of placeholder refer Get or GetAll
// it returns number of rows modified and error
func (databaseRepo zkPostgresRepo) Delete(tx *sql.Tx, query string, param []any) (int, error) {
	return databaseRepo.modifyTable(tx, query, param)
}

// Update This method takes a transaction, sql query, slice of params where each value in the param provides values to
// placeholders in the query. For more details of placeholder refer Get or GetAll
// it returns number of rows modified and error
func (databaseRepo zkPostgresRepo) Update(tx *sql.Tx, stmt string, param []any) (int, error) {
	return databaseRepo.modifyTable(tx, stmt, param)
}

// internal method used by update and delete
func (databaseRepo zkPostgresRepo) modifyTable(tx *sql.Tx, stmt string, param []any) (int, error) {
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

// Insert This method takes a db handle, sql query and DbArgs
// DbArgs is an interface which defines a method GetAllColumns which returns a slice of struct fields corresponding to columns
// in db table, The values from this slice is then inserted into the table.
func (databaseRepo zkPostgresRepo) Insert(db *sql.DB, query string, data interfaces.DbArgs) (sql.Result, error) {
	preparedStmt, err := db.Prepare(query)
	if err != nil {
		zkLogger.Error(LogTag, "Error preparing insert statement:", err)
		return nil, err
	}
	defer preparedStmt.Close()

	c := data.GetAllColumns()
	result, err := preparedStmt.Exec(c)
	if err != nil {
		zkLogger.Error(LogTag, "Error executing insert statement:", err)
		return nil, err
	}

	return result, nil
}

// BulkInsert This method takes a transaction, tableName, list of columns and corresponding column values in []DbArgs
// For more details on DbArgs read Insert
// Here we are using copyIn which is a very fast operation for bulk Insert
func (databaseRepo zkPostgresRepo) BulkInsert(tx *sql.Tx, tableName string, columns []string, data []interfaces.DbArgs) error {
	stmt, err := tx.Prepare(pq.CopyIn(tableName, columns...))
	if err != nil {
		zkLogger.Error(LogTag, "Error preparing insert statement:", err)
		return err
	}

	for _, d := range data {
		c := d.GetAllColumns()
		_, err := stmt.Exec(c...)
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

// InsertInTransaction This method takes a db handle, sql query and DbArgs
// DbArgs is an interface which defines a method GetAllColumns which returns a slice of struct fields corresponding to columns
// in db table, The values from this slice is then inserted into the table.
func (databaseRepo zkPostgresRepo) InsertInTransaction(tx *sql.Tx, stmt string, data interfaces.DbArgs) (sql.Result, error) {
	preparedStmt, err := tx.Prepare(stmt)
	if err != nil {
		zkLogger.Error(LogTag, "Error preparing insert statement:", err)
		return nil, err
	}
	defer preparedStmt.Close()

	c := data.GetAllColumns()
	Result, err := preparedStmt.Exec(c...)
	if err != nil {
		zkLogger.Error(LogTag, "Error executing insert:", err)
		return nil, err
	}

	return Result, nil
}

// BulkUpsert This method takes a transaction, sql query and a slice DbArgs
// For more details on DbArgs read Insert
// Here we are using Upsert command which gives additional control when inserts operation fail due to conflicts.
// Example Upsert Query: INSERT INTO TableA (colA, colB) VALUES ($1, $2) ON CONFLICT (colA) DO NOTHING
// The above query tries to insert some value to TableA but on conflict, it does nothing, you can also specify what to do
// incase of a conflict.
func (databaseRepo zkPostgresRepo) BulkUpsert(tx *sql.Tx, stmt string, data []interfaces.DbArgs) error {
	preparedStmt, err := tx.Prepare(stmt)
	if err != nil {
		zkLogger.Error(LogTag, "Error preparing insert statement:", err)
		return err
	}

	defer preparedStmt.Close()

	for _, d := range data {
		c := d.GetAllColumns()
		_, err = preparedStmt.Exec(c...)
		if err != nil {
			_ = tx.Rollback()
			zkLogger.Error(LogTag, "failed to perform bulk upsert: ", err)
			return err
		}
	}

	return nil
}
