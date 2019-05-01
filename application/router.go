package application

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (a *App) Router() *chi.Mux {
	r := chi.NewRouter()

	// More information on Chi middleware can be found at https://github.com/go-chi/chi#middlewares
	r.Use(middleware.Heartbeat("_health"))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(StripSlashes)

	r.Route("/project", func(r chi.Router) {
		r.Get("/", a.ListProjects)
		r.Get("/{name}", a.GetProject)
		r.Get("/{name}/{instance}", a.GetProjectInstance)

		r.Post("/", a.CreateProject)
		r.Put("/{name}", a.UpdateProject)
		r.Delete("/{name}", a.DeleteProject)

		r.Post("/{name}", a.CreateProjectInstance)
		r.Put("/{name}/{instance}", a.UpdateProjectInstance)
		r.Delete("/{name}/{instance}", a.DeleteProjectInstance)
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/", a.ListUsers)
		r.Post("/", a.CreateUser)
		r.Put("/", a.UpdateUser)
		r.Get("/{name}", a.GetUser)
		r.Delete("/{name}", a.DeleteUser)
	})

	// Execute instance
	r.Route("/e", func(r chi.Router) {
		r.Get("/{project}", a.ExecuteInstance)
		r.Get("/{project}/{alias}", a.ExecuteInstance)
		r.Post("/{project}", a.ExecuteInstance)
		r.Post("/{project}/{alias}", a.ExecuteInstance)
	})

	return r
}
