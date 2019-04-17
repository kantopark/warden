package server

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/spf13/viper"

	"warden/store"
)

type Server struct {
	Docker *client.Client
	Store  *store.Store
}

func NewServer() *Server {
	docker, err := client.NewEnvClient()
	if err != nil {
		panic(fmt.Errorf("error creating docker client: %s\n", err))
	}

	dialect := viper.GetString("store.dialect")
	dsn := viper.GetString("store.dsn")
	_store, err := store.NewStore(dialect, dsn)
	if err != nil {
		panic(fmt.Errorf("error creating store: %s\n", err))
	}

	return &Server{
		Docker: docker,
		Store:  _store,
	}
}
