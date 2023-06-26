package test

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zerok-ai/zk-utils-go/common"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB/postgres"
	"github.com/zerok-ai/zk-utils-go/storage/sqlDB/postgres/config"
	"log"
	"testing"

	// The "testify/suite" package is used to make the test suite
	"github.com/stretchr/testify/suite"
)

const (
	InsertTraceQuery              = "INSERT INTO trace (scenario_id, scenario_version, trace_id) VALUES ($1, $2, $3)"
	UpdateTraceQuery              = "UPDATE trace SET scenario_version=$1 WHERE trace_id=$2"
	DeleteTraceQuery              = "DELETE FROM trace WHERE scenario_version=$1"
	GetAllTraceQuery              = "SELECT scenario_id, scenario_version, trace_id FROM trace"
	GetAllTraceQueryWithCondition = "SELECT scenario_id, scenario_version, trace_id FROM trace WHERE scenario_version=$1 AND trace_id=$2"
	GetTraceQuery                 = "SELECT scenario_id, scenario_version, trace_id FROM trace WHERE scenario_version=$1"
)

var LogTag = "zkpostgres_test"

type StoreSuite struct {
	suite.Suite
	/*
		The suite is defined as a struct, with the store and db as its
		attributes. Any variables that are to be shared between tests in a
		suite should be stored as attributes of the suite instance
	*/
	dbRepo sqlDB.DatabaseRepo
}

func (s *StoreSuite) SetupSuite() {
	/*
		The database connection is opened in the setup, and
		stored as an instance variable,
		as is the higher level `store`, that wraps the `db`
	*/
	c := config.PostgresConfig{
		Host:                           "localhost",
		Port:                           5432,
		User:                           "pl",
		Password:                       "pl",
		Dbname:                         "pl",
		MaxConnections:                 0,
		MaxIdleConnections:             0,
		ConnectionMaxLifetimeInMinutes: 0,
	}

	r, err := zkpostgres.NewZkPostgresRepo(c)
	if err != nil {
		log.Fatal("unable to connect to db")
	}
	s.dbRepo = r
}

func (s *StoreSuite) SetupTest() {
	/*
		We delete all entries from the table before each test runs, to ensure a
		consistent state before our tests run. In more complex applications, this
		is sometimes achieved in the form of migrations
	*/
	tx, err := s.dbRepo.CreateTransaction()
	if err != nil {
		s.T().Fatal(err)
	}

	stmt, _ := common.GetStmtRawQuery(tx, "DELETE FROM trace")

	_, err = s.dbRepo.Delete(stmt, nil)
	if err != nil {
		tx.Rollback()
		s.T().Fatal(err)
	}
	tx.Commit()
}

func (s *StoreSuite) TearDownSuite() {
	// Close the connection after all tests in the suite finish
	err := s.dbRepo.Close()
	if err != nil {
		return
	}
}

// This is the actual "test" as seen by Go, which runs the tests defined below
func TestStoreSuite(t *testing.T) {
	s := new(StoreSuite)
	suite.Run(t, s)
}

func (s *StoreSuite) TestInsertTrace() {
	tx, err := s.dbRepo.CreateTransaction()
	if err != nil {
		s.T().Fatal(err)
	}

	stmt, _ := common.GetStmtRawQuery(tx, InsertTraceQuery)
	t := getTrace(1)

	// Create a bird through the store `CreateBird` method
	insert, err := s.dbRepo.Insert(stmt, t)
	if err != nil {
		s.T().Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		s.T().Fatal(err)
	}

	c, err := insert.RowsAffected()
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), 1, int(c))

	// Query the database for the entry we just created
	rows, err, closeRow := s.dbRepo.GetAll(GetAllTraceQuery, nil)
	defer closeRow()

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), rows)

	var responseArr []Trace
	for rows.Next() {
		var rawData Trace
		err := rows.Scan(&rawData.ScenarioId, &rawData.ScenarioVersion, &rawData.TraceId)
		if err != nil {
			log.Fatal(err)
		}

		responseArr = append(responseArr, rawData)
	}

	assert.Equal(s.T(), 1, len(responseArr))
	assert.Equal(s.T(), t.TraceId, responseArr[0].TraceId)
	assert.Equal(s.T(), t.ScenarioId, responseArr[0].ScenarioId)
	assert.Equal(s.T(), t.ScenarioVersion, responseArr[0].ScenarioVersion)
}

func (s *StoreSuite) TestGetAllTrace() {
	tx, err := s.dbRepo.CreateTransaction()
	if err != nil {
		s.T().Fatal(err)
	}

	expected1 := getTrace(1)
	expected2 := getTrace(2)

	stmt, _ := common.GetStmtRawQuery(tx, InsertTraceQuery)
	insert1, err := s.dbRepo.Insert(stmt, expected1)

	stmt, _ = common.GetStmtRawQuery(tx, InsertTraceQuery)
	insert2, err := s.dbRepo.Insert(stmt, expected2)
	if err != nil {
		s.T().Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		s.T().Fatal(err)
	}

	c1, err := insert1.RowsAffected()
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), 1, int(c1))

	c2, err := insert2.RowsAffected()
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), 1, int(c2))

	// GET ALL WITHOUT CONDITION
	rows, err, closeRow := s.dbRepo.GetAll(GetAllTraceQuery, nil)
	defer closeRow()

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), rows)

	var responseArr []Trace
	for rows.Next() {
		var rawData Trace
		err := rows.Scan(&rawData.ScenarioId, &rawData.ScenarioVersion, &rawData.TraceId)
		if err != nil {
			log.Fatal(err)
		}

		responseArr = append(responseArr, rawData)
	}

	actual1 := responseArr[0]
	actual2 := responseArr[1]
	assert.Equal(s.T(), 2, len(responseArr))
	assert.Equal(s.T(), expected1.TraceId, actual1.TraceId)
	assert.Equal(s.T(), expected1.ScenarioId, actual1.ScenarioId)
	assert.Equal(s.T(), expected1.ScenarioVersion, actual1.ScenarioVersion)

	assert.Equal(s.T(), expected2.TraceId, actual2.TraceId)
	assert.Equal(s.T(), expected2.ScenarioId, actual2.ScenarioId)
	assert.Equal(s.T(), expected2.ScenarioVersion, actual2.ScenarioVersion)

	// GET ALL WITH CONDITION
	rows, err, closeRow = s.dbRepo.GetAll(GetAllTraceQueryWithCondition, []any{"sv:2", "tId:2"})
	defer closeRow()

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), rows)
	responseArr = responseArr[:0]

	for rows.Next() {
		var rawData Trace
		err := rows.Scan(&rawData.ScenarioId, &rawData.ScenarioVersion, &rawData.TraceId)
		if err != nil {
			log.Fatal(err)
		}

		responseArr = append(responseArr, rawData)
	}

	actual1 = responseArr[0]
	assert.Equal(s.T(), 1, len(responseArr))
	assert.Equal(s.T(), expected2.TraceId, actual1.TraceId)
	assert.Equal(s.T(), expected2.ScenarioId, actual1.ScenarioId)
	assert.Equal(s.T(), expected2.ScenarioVersion, actual1.ScenarioVersion)
}

func (s *StoreSuite) TestGetTrace() {
	tx, err := s.dbRepo.CreateTransaction()
	if err != nil {
		s.T().Fatal(err)
	}

	expected1 := getTrace(1)
	expected2 := getTrace(2)

	stmt, _ := common.GetStmtRawQuery(tx, InsertTraceQuery)
	insert1, err := s.dbRepo.Insert(stmt, expected1)

	stmt, _ = common.GetStmtRawQuery(tx, InsertTraceQuery)
	insert2, err := s.dbRepo.Insert(stmt, expected2)
	if err != nil {
		s.T().Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		s.T().Fatal(err)
	}

	c1, err := insert1.RowsAffected()
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), 1, int(c1))

	c2, err := insert2.RowsAffected()
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), 1, int(c2))

	var result Trace
	// GET WITHOUT CONDITION
	err = s.dbRepo.Get(GetTraceQuery, []any{"sv:1"}, []any{&result.ScenarioId, &result.ScenarioVersion, &result.TraceId})

	assert.Nil(s.T(), err)

	assert.Equal(s.T(), expected1.TraceId, result.TraceId)
	assert.Equal(s.T(), expected1.ScenarioId, result.ScenarioId)
	assert.Equal(s.T(), expected1.ScenarioVersion, result.ScenarioVersion)
}

func (s *StoreSuite) TestUpdateTrace() {
	tx, err := s.dbRepo.CreateTransaction()
	if err != nil {
		s.T().Fatal(err)
	}

	stmt, _ := common.GetStmtRawQuery(tx, InsertTraceQuery)
	t := getTrace(1)

	// Create a bird through the store `CreateBird` method
	insert, err := s.dbRepo.Insert(stmt, t)
	if err != nil {
		s.T().Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		s.T().Fatal(err)
	}

	c, err := insert.RowsAffected()
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), 1, int(c))

	// Query the database for the entry we just created
	rows, err, closeRow := s.dbRepo.GetAll(GetAllTraceQuery, nil)
	defer closeRow()

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), rows)

	var responseArr []Trace
	for rows.Next() {
		var rawData Trace
		err := rows.Scan(&rawData.ScenarioId, &rawData.ScenarioVersion, &rawData.TraceId)
		if err != nil {
			log.Fatal(err)
		}

		responseArr = append(responseArr, rawData)
	}

	assert.Equal(s.T(), 1, len(responseArr))
	assert.Equal(s.T(), t.TraceId, responseArr[0].TraceId)
	assert.Equal(s.T(), t.ScenarioId, responseArr[0].ScenarioId)
	assert.Equal(s.T(), t.ScenarioVersion, responseArr[0].ScenarioVersion)

	//************************** NOW UPDATE THE RECORD
	tx, err = s.dbRepo.CreateTransaction()
	if err != nil {
		s.T().Fatal(err)
	}

	stmt, _ = common.GetStmtRawQuery(tx, UpdateTraceQuery)

	// Create a bird through the store `CreateBird` method
	newScenarioVersion := "new_scenario_id"
	traceId := "tId:1"
	update, err := s.dbRepo.Delete(stmt, []any{newScenarioVersion, traceId})
	if err != nil {
		s.T().Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		s.T().Fatal(err)
	}

	c, err = update.RowsAffected()
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), 1, int(c))

	var result Trace
	// GET WITHOUT CONDITION
	err = s.dbRepo.Get(GetTraceQuery, []any{newScenarioVersion}, []any{&result.ScenarioId, &result.ScenarioVersion, &result.TraceId})

	assert.Nil(s.T(), err)

	assert.Equal(s.T(), traceId, result.TraceId)
	assert.Equal(s.T(), t.ScenarioId, result.ScenarioId)
	assert.Equal(s.T(), newScenarioVersion, result.ScenarioVersion)
}

func (s *StoreSuite) TestDeleteTrace() {
	tx, err := s.dbRepo.CreateTransaction()
	if err != nil {
		s.T().Fatal(err)
	}

	stmt, _ := common.GetStmtRawQuery(tx, InsertTraceQuery)
	t := getTrace(1)

	// Create a bird through the store `CreateBird` method
	insert, err := s.dbRepo.Insert(stmt, t)
	if err != nil {
		s.T().Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		s.T().Fatal(err)
	}

	c, err := insert.RowsAffected()
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), 1, int(c))

	// Query the database for the entry we just created
	rows, err, closeRow := s.dbRepo.GetAll(GetAllTraceQuery, nil)
	defer closeRow()

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), rows)

	var responseArr []Trace
	for rows.Next() {
		var rawData Trace
		err := rows.Scan(&rawData.ScenarioId, &rawData.ScenarioVersion, &rawData.TraceId)
		if err != nil {
			log.Fatal(err)
		}

		responseArr = append(responseArr, rawData)
	}

	scenarioVersion := "sv:1"

	assert.Equal(s.T(), 1, len(responseArr))
	assert.Equal(s.T(), t.TraceId, responseArr[0].TraceId)
	assert.Equal(s.T(), t.ScenarioId, responseArr[0].ScenarioId)
	assert.Equal(s.T(), scenarioVersion, responseArr[0].ScenarioVersion)

	//************************** NOW UPDATE THE RECORD
	tx, err = s.dbRepo.CreateTransaction()
	if err != nil {
		s.T().Fatal(err)
	}

	stmt, _ = common.GetStmtRawQuery(tx, DeleteTraceQuery)

	// Create a bird through the store `CreateBird` method
	update, err := s.dbRepo.Delete(stmt, []any{scenarioVersion})
	if err != nil {
		s.T().Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		s.T().Fatal(err)
	}

	c, err = update.RowsAffected()
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), 1, int(c))

	var result Trace
	// GET WITHOUT CONDITION
	err = s.dbRepo.Get(GetTraceQuery, []any{scenarioVersion}, []any{&result.ScenarioId, &result.ScenarioVersion, &result.TraceId})

	assert.Equal(s.T(), sql.ErrNoRows, err)

	assert.Empty(s.T(), result.TraceId)
	assert.Empty(s.T(), result.ScenarioId)
	assert.Empty(s.T(), result.ScenarioVersion)
}

type Trace struct {
	ScenarioId      int    `json:"scenario_id"`
	ScenarioVersion string `json:"scenario_version"`
	TraceId         string `json:"trace_id"`
}

func getTrace(i int) Trace {
	return Trace{
		ScenarioId:      i,
		ScenarioVersion: fmt.Sprintf("sv:%d", i),
		TraceId:         fmt.Sprintf("tId:%d", i),
	}
}

func (t Trace) GetAllColumns() []any {
	return []any{t.ScenarioId, t.ScenarioVersion, t.TraceId}
}
