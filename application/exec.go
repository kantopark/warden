package application

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Executes the function specified
func (a *App) ExecuteInstance(w http.ResponseWriter, r *http.Request) {
	project := chi.URLParam(r, "project")
	alias := chi.URLParam(r, "alias")
	if alias == "" {
		alias = "latest"
	}

	panic(project)
}
