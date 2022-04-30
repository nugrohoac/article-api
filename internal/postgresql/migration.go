package postgresql

import (
	"log"
	"path"
	"runtime"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Migration is struct for migration purpose
type Migration struct {
	Migrate *migrate.Migrate
}

// Up is method for migrate up database
func (m *Migration) Up() ([]error, bool) {
	err := m.Migrate.Up()
	if err != nil {
		return []error{err}, false
	}

	return []error{}, true
}

// Down is method for migrate down database
func (m *Migration) Down() ([]error, bool) {
	err := m.Migrate.Down()
	if err != nil {
		return []error{err}, false
	}

	return []error{}, true
}

// RunMigration is function to run the database migration (up and down)
func RunMigration(dbURI string) (*Migration, error) {
	_, filename, _, _ := runtime.Caller(0)

	migrationPath := path.Join(path.Dir(filename), "migrations")

	var dataPath []string
	dataPath = append(dataPath, "file://")
	dataPath = append(dataPath, migrationPath)

	pathToMigrate := strings.Join(dataPath, "")

	m, err := migrate.New(pathToMigrate, dbURI)
	if err != nil {
		log.Fatal(err)
	}

	return &Migration{Migrate: m}, err
}
