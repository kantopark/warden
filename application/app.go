package application

import (
	"log"

	"warden/docker"
	"warden/store"
)

type App struct {
	dck *docker.Client
	db  *store.Store
}

// Creates a new App object
func NewApp() *App {
	_dck, err := docker.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	_db, err := store.NewStore()
	if err != nil {
		log.Fatalf("error creating store: %s\n", err)
	}

	app := &App{
		dck: _dck,
		db:  _db,
	}

	return app
}

func (a *App) Close() (err error) {
	err = a.db.Close()
	err = a.dck.Close()
	return err
}
