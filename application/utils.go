package application

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// The base error response body
func errorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	msg := ""
	if err != nil {
		msg = err.Error()
	}

	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{Error: msg})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Returns a Bad Request response
func badRequest(w http.ResponseWriter, err error) {
	errorResponse(w, err, http.StatusBadRequest)
}

// Returns interface object as json.
func jsonify(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		internalServerError(w, err)
	}
}

// Returns a Status Okay response
func ok(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(struct {
		Result string `json:"result"`
	}{"Ok"}); err != nil {
		internalServerError(w, err)
	}
}

// Utility function for marshalling the request's body into the interface
func parseJson(body io.ReadCloser, v interface{}) error {
	defer body.Close()
	if err := json.NewDecoder(body).Decode(v); err != nil {
		return errors.Wrap(err, "error decoding request payload")
	}
	return nil
}

// Returns an Internal Server Error response
func internalServerError(w http.ResponseWriter, err error) {
	errorResponse(w, err, http.StatusInternalServerError)
}
