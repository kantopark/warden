package router

import "net/http"

// Post request. Creates a new function in the system. JSON payload
// should specify creation parameters
func (a *app) CreateFunction(w http.ResponseWriter, r *http.Request) {

}

// Post request. Appends a new RunInfo instance to function.
func (a *app) CreateRunInfo(w http.ResponseWriter, r *http.Request) {

}

// Delete request. Removes function and all RunInfo associated with it
func (a *app) DeleteFunction(w http.ResponseWriter, r *http.Request) {

}

// Delete request. Removes RunInfo from function. Cannot remove default
// RunInfo (alias == "latest" or empty string).
func (a *app) DeleteRunInfo(w http.ResponseWriter, r *http.Request) {

}

// Get request. List all functions
func (a *app) ListFunctions(w http.ResponseWriter, r *http.Request) {

}

// Get request. List all RunInfo related to function
func (a *app) ListRunInfo(w http.ResponseWriter, r *http.Request) {

}

// Get request. Reads information related to a specific function
func (a *app) ReadFunction(w http.ResponseWriter, r *http.Request) {

}

// Get request. Reads RunInfo data related to a specific function
func (a *app) ReadRunInfo(w http.ResponseWriter, r *http.Request) {

}

// Put request. Updates function with JSON payload specified
func (a *app) UpdateFunction(w http.ResponseWriter, r *http.Request) {

}

// Put Request. Updates a RunInfo associated with the function with
// the JSON payload
func (a *app) UpdateRunInfo(w http.ResponseWriter, r *http.Request) {

}
