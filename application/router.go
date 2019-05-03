package application

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
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

	// Protected paths
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwtToken))
		r.Use(JWTAuthenticate)

		r.Route("/project", func(r chi.Router) {
			r.Get("/", a.ListProjects)
			r.Get("/{name}", a.GetProject)
			r.Post("/", a.CreateProject)
			r.Put("/", a.UpdateProject)
			r.Delete("/{name}", a.DeleteProject)
		})

		r.Route("/project-instance", func(r chi.Router) {
			r.Post("/", a.CreateProjectInstance)
			r.Put("/", a.UpdateProjectInstance)
			r.Delete("/{name}/{id}", a.DeleteProjectInstance)
		})

		r.Route("/user", func(r chi.Router) {
			r.Get("/", a.ListUsers)
			r.Put("/", a.UpdateUser)
			r.Get("/{name}", a.GetUser)
			r.Delete("/{name}", a.DeleteUser)
		})
	})

	// public routes
	r.Group(func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", a.Login)
			r.Post("/signup", a.Signup)
		})
		// Execute instance
		r.Route("/e", func(r chi.Router) {
			r.Get("/{project}", a.ExecuteInstance)
			r.Get("/{project}/{alias}", a.ExecuteInstance)
			r.Post("/{project}", a.ExecuteInstance)
			r.Post("/{project}/{alias}", a.ExecuteInstance)
		})
	})

	return r
}
