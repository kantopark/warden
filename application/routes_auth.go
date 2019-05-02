package application

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"warden/store"
)

// Post request. If successful, will return a JWT token
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	var u store.UserBody
	if err := parseJson(r.Body, &u); err != nil {
		internalServerError(w, err)
		return
	}
	user, err := a.db.UserLogin(u, true)
	if err != nil {
		internalServerError(w, err)
	}

	token, err := createToken(user.Username, user.Email)
	if err != nil {
		internalServerError(w, errors.Wrap(err, "error creating jwt token"))
		return
	}

	w.Header().Set("authorization", fmt.Sprintf("Bearer %s", token))
	jsonify(w, struct {
		Token string `json:"token"`
	}{token})
}

// Post request. Creates a new user. If successful, will return a JWT token
func (a *App) Signup(w http.ResponseWriter, r *http.Request) {
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

	token, err := createToken(user.Username, user.Email)
	if err != nil {
		internalServerError(w, errors.Wrap(err, "error creating jwt token"))
		return
	}

	w.Header().Set("authorization", fmt.Sprintf("Bearer %s", token))
	jsonify(w, struct {
		Token string `json:"token"`
	}{token})
}
