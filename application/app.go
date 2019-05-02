package application

import (
	"warden/docker"
	"warden/store"
)

type App struct {
	dck *docker.Client
	db  *store.Store
}

// Creates a new App object
func NewApp() *App {
	authInit()

	_dck, err := docker.NewClient()
	fatalIfError(err)

	_db, err := store.NewStore()
	fatalIfError(err)

	app := &App{
		dck: _dck,
		db:  _db,
	}

	return app
}

func (a *App) Close() {
	err := a.db.Close()
	fatalIfError(err)

	// Ignoring docker close
	//err = a.dck.Close()
	//fatalIfError(err)
}
