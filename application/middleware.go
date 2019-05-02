package application

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

const userContextKey = "user_context_key"

type userCtx struct {
	Email    string
	Exp      time.Time
	Username string
}

// StripSlashes is a middleware that will match request paths with a trailing
// slash, strip it from the path and continue routing through the mux, if a route
// matches, then it will serve the handler.
func StripSlashes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var path string
		rctx := chi.RouteContext(r.Context())
		if rctx.RoutePath != "" {
			path = rctx.RoutePath
		} else {
			path = r.URL.Path
		}
		for len(path) > 1 && path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
		rctx.RoutePath = path
		next.ServeHTTP(w, r)
	})
}

func JWTAuthenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, err := jwtToken.Decode(bearerToken[1])
				if err != nil {
					internalServerError(w, errors.Wrap(err, "error decoding JWT token"))
					return
				}

				// Casting claims into userCtx object. Downstream handlers will be able to retrieve value
				// by calling `r.Context().Value(userContextKey).(userCtx)`
				claims := token.Claims.(jwt.MapClaims)
				u := userCtx{
					Email:    claims["email"].(string),
					Exp:      time.Unix(int64(claims["exp"].(float64)), 0),
					Username: claims["username"].(string),
				}

				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, u)))
				return
			}
		}
		badRequest(w, errors.New("Invalid authorization header"))
	})
}
