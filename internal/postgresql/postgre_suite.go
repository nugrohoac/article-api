package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// TestSuite is struct for test suite
type TestSuite struct {
	suite.Suite
	DSN    string
	M      *Migration
	DBConn *sql.DB
}

// default value for CI/CD
var (
	// postgresSQL
	driver   = "postgres"
	host     = "localhost"
	dbname   = "kumparan_assessment_test"
	sslMode  = "disable"
	userName = "postgres"
	password = "password123"

	searchPath = "public"
	port       = "5432"
)

// SetupSuite is method for set up the test suite
func (s *TestSuite) SetupSuite() {
	var err error
	s.DSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s", userName, password, host, port, dbname, sslMode, searchPath)

	s.DBConn, err = sql.Open(driver, s.DSN)
	require.NoError(s.T(), err)

	err = s.DBConn.Ping()
	require.NoError(s.T(), err)

	s.M, err = RunMigration(s.DSN)
	require.NoError(s.T(), err)

	//s.cleanIfDirtyDatabase()

	logrus.Info("Starting to Migrate Up Data")
	errUp, okUp := s.M.Up()
	for _, element := range errUp {
		require.NoError(s.T(), element)
	}

	require.True(s.T(), okUp)
	require.Len(s.T(), errUp, 0)
}

// TearDownSuite is method which will be run when the test suite is done
func (s *TestSuite) TearDownSuite() {
	logrus.Info("Starting to Migrate Down Data")
	err, ok := s.M.Down()
	require.True(s.T(), ok)
	require.Len(s.T(), err, 0)

	errClose := s.DBConn.Close()
	require.NoError(s.T(), errClose)
}

// TearDownTest is called when starting migrate down
func (s *TestSuite) TearDownTest() {
	query := `SELECT TABLE_NAME FROM information_schema.tables WHERE table_schema='` + searchPath + `'`
	rows, err := s.DBConn.Query(query)
	require.NoError(s.T(), err)

	for rows.Next() {
		var tableName string

		if errScan := rows.Scan(&tableName); errScan != nil {
			logrus.Error("error_message_scan : ", errScan)
		}

		if tableName == "schema_migrations" {
			continue
		}

		if tableName == "user" {
			tableName = "\"" + tableName + "\""
		}

		queryTruncate := "TRUNCATE TABLE " + tableName
		_, err = s.DBConn.Exec(queryTruncate)
		require.NoError(s.T(), err)
	}

	err = rows.Close()
	require.NoError(s.T(), err)
}
