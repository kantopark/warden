package application

import (
	"net/http"

	"github.com/go-chi/chi"

	"warden/store"
)

// Delete request. Removes a user (group)
func (a *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if err := a.db.UserDelete(name); err != nil {
		internalServerError(w, err)
		return
	}
	ok(w)
}

// Get request. List all user (groups)
func (a *App) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := a.db.UserList(true)
	if err != nil {
		internalServerError(w, err)
		return
	}
	jsonify(w, users)
}

// Get request. Reads information about user (group)
func (a *App) GetUser(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	u, err := a.db.UserGet(name, true)
	if err != nil {
		internalServerError(w, err)
		return
	}
	jsonify(w, u)
}

// Put request. Updates information about user.
func (a *App) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// TODO allow admin to overwrite this
	// TODO add password reset
	var u store.UserBody

	if err := parseJson(r.Body, &u); err != nil {
		internalServerError(w, err)
		return
	}

	user, err := a.db.UserUpdate(u)
	if err != nil {
		internalServerError(w, err)
		return
	}

	jsonify(w, user)
}
