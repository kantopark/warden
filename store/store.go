package store

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"

	"warden/store/model"
)

const (
	_INSTANCE = "Instance"
	_PROJECT  = "Project"
	_USER     = "User"
)

// A store used to carry information on the functions.
type Store struct {
	db *gorm.DB
}

// Creates a new store object. The store is used to hold information on how to run
// functions
func NewStore(dialect, dsn string) (*Store, error) {
	db, err := gorm.Open(dialect, dsn)

	if err != nil {
		return nil, errors.Wrapf(err, "error creating database connection. Dialect: '%s'. DSN: '%s'", dialect, dsn)
	}

	s := &Store{db}
	s.registerModels()

	return s, nil
}

// Registers all models and migrates database to the latest version
func (s *Store) registerModels() {
	s.CreateTableIfNotExists(&model.User{})
	s.CreateTableIfNotExists(&model.Project{})
	s.CreateTableIfNotExists(&model.Instance{})
}

// Creates table if it doesn't exist. Else migrates the table to the latest state.
// Migration will only add missing fields for the given model and won't
// delete/change current data
func (s *Store) CreateTableIfNotExists(table interface{}) {
	if !s.db.HasTable(table) {
		s.db.CreateTable(table)
	} else {
		s.db.AutoMigrate(table)
	}
}

func (s *Store) Close() error {
	return s.db.Close()
}
