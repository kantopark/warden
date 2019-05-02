package store

import (
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"warden/store/model"
)

const (
	_INSTANCES = "Instances"
	_PROJECTS  = "Projects"
	_OWNERS    = "Owners"
)

// A store used to carry information on the functions.
type Store struct {
	db *gorm.DB
}

var store *Store
var storeError error
var once sync.Once

// Returns a Store object. The store is used to hold information on how to run functions
// The store is implemented as a singleton. Subsequent to NewStore will only return the
// first created instance of the store.
func NewStore() (*Store, error) {
	once.Do(func() {
		dialect := viper.GetString("store.dialect")
		dsn := viper.GetString("store.dsn")
		db, err := gorm.Open(dialect, dsn)
		if err != nil {
			storeError = errors.Wrapf(err, "error creating database connection. Dialect: '%s'. DSN: '%s'", dialect, dsn)
			return
		}

		db.LogMode(viper.GetBool("store.log_mode"))

		store = &Store{db}
		store.registerModels()
	})

	return store, storeError
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
