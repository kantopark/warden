package store

import (
	"log"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"warden/config"
)

var (
	S        *Store
	username = "daniel"
)

func init() {
	var err error

	config.ReadConfig()
	viper.Set("store.dsn", "file:memdb1?mode=memory&cache=shared")
	viper.Set("store.dialect", "sqlite3")

	S, err = NewStore()
	if err != nil {
		log.Fatalln(err)
	}

	user, err := S.UserCreate(UserBody{
		Email:    "daniel.bok@outlook.com",
		Username: username,
		Password: "password",
	})
	if err != nil {
		log.Fatalln(err)
	}

	project, err := S.ProjectCreate(
		"https://github.com/kantopark-tpl/python-simple",
		"python-test",
		"A simple description",
		*user)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = S.InstanceCreate("95bfc3515452bfafeb2e04f948ac26d1e2a871c8", "test", project.Name)
	if err != nil {
		log.Fatalln(err)
	}
}

func TestNewStore(t *testing.T) {
	_, err := NewStore()
	assert.Nil(t, err)
}
