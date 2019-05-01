package application

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"warden/store"
)

// Post request. Creates a new user (group)
func (a *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u store.UserBody
	if err := parseJson(r.Body, &u); err != nil {
		internalServerError(w, err)
		return
	}

	minLen := viper.GetInt("auth.pw_len")
	if len(u.Password) < minLen {
		badRequest(w, errors.Errorf("Password must be >= %d characters", minLen))
		return
	}

	user, err := a.db.UserCreate(u)
	if err != nil {
		internalServerError(w, err)
		return
	}

	jsonify(w, user)
}

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

	user, err := a.db.UserUpdatePassword(u)
	if err != nil {
		internalServerError(w, err)
		return
	}

	jsonify(w, user)
}
