package router

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"

	"warden/store"
)

type app struct {
	Docker *client.Client
	Store  *store.Store
}

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	// More information on Chi middleware can be found at https://github.com/go-chi/chi#middlewares
	r.Use(middleware.Heartbeat("_health"))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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

	a := &app{
		Docker: docker,
		Store:  _store,
	}

	registerRoutes(r, a)

	return r
}

func registerRoutes(r *chi.Mux, a *app) {
	r.Route("function", func(r chi.Router) {
		r.Get("/", a.ListFunctions)
		r.Get("/{name}", a.ReadFunction)

		r.Post("/", a.CreateFunction)
		r.Put("/", a.UpdateFunction)
		r.Delete("/", a.DeleteFunction)
	})

	r.Route("runinfo", func(r chi.Router) {
		r.Get("/{name}", a.ListRunInfo)
		r.Post("/", a.CreateRunInfo)
		r.Put("/", a.UpdateRunInfo)
		r.Delete("/", a.DeleteRunInfo)
	})

	r.Route("exec", func(r chi.Router) {
		r.Get("/{name}", a.ExecuteFunctionGet)
		r.Get("/{name}/{alias}", a.ExecuteFunctionGet)

		r.Post("/{name}", a.ExecuteFunctionPost)
		r.Post("/{name}/alias", a.ExecuteFunctionPost)
	})
}
