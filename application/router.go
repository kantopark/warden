package application

import (
	"log"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"

	"warden/docker"
	"warden/store"
)

type App struct {
	cli *docker.Client
	db  *store.Store
}

// Creates a new App object
func NewApp() *App {
	_cli, err := docker.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	dialect := viper.GetString("store.dialect")
	dsn := viper.GetString("store.dsn")
	_db, err := store.NewStore(dialect, dsn)
	if err != nil {
		log.Fatalf("error creating store: %s\n", err)
	}

	app := &App{
		cli: _cli,
		db:  _db,
	}

	return app
}

func (a *App) Close() (err error) {
	err = a.db.Close()
	err = a.cli.Close()
	return err
}

func (a *App) Router() *chi.Mux {
	r := chi.NewRouter()

	// More information on Chi middleware can be found at https://github.com/go-chi/chi#middlewares
	r.Use(middleware.Heartbeat("_health"))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/function", func(r chi.Router) {
		r.Get("/", a.ListFunctions)
		r.Get("/{name}", a.ReadFunction)

		r.Post("/", a.CreateFunction)
		r.Put("/", a.UpdateFunction)
		r.Delete("/", a.DeleteFunction)
	})

	r.Route("/runinfo", func(r chi.Router) {
		r.Get("/{name}", a.ListRunInfo)
		r.Post("/", a.CreateRunInfo)
		r.Put("/", a.UpdateRunInfo)
		r.Delete("/", a.DeleteRunInfo)
	})

	r.Route("/exec", func(r chi.Router) {
		r.Get("/{name}", a.ExecuteFunctionGet)
		r.Get("/{name}/{alias}", a.ExecuteFunctionGet)

		r.Post("/{name}", a.ExecuteFunctionPost)
		r.Post("/{name}/alias", a.ExecuteFunctionPost)
	})

	return r
}
