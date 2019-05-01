package application

import "net/http"

// Post request. Creates a new project in the system. JSON payload
// should specify creation parameters
func (a *App) CreateProject(w http.ResponseWriter, r *http.Request) {

}

// Post request. Appends a new Instances to Project.
func (a *App) CreateProjectInstance(w http.ResponseWriter, r *http.Request) {

}

// Post request. Creates a new user (group)
func (a *App) CreateUser(w http.ResponseWriter, r *http.Request) {

}

// Delete request. Removes project and all Instances associated with it
func (a *App) DeleteProject(w http.ResponseWriter, r *http.Request) {

}

// Delete request. Removes an instance from the project. Project must however be left with
// at least the latest instance.
func (a *App) DeleteProjectInstance(w http.ResponseWriter, r *http.Request) {

}

// Delete request. Removes a user (group)
func (a *App) DeleteUser(w http.ResponseWriter, r *http.Request) {

}

// Get request. List all project
func (a *App) ListProjects(w http.ResponseWriter, r *http.Request) {

}

// Get request. List all user (groups)
func (a *App) ListUsers(w http.ResponseWriter, r *http.Request) {

}

// Get request. Reads information related to a specific project
func (a *App) GetProject(w http.ResponseWriter, r *http.Request) {

}

// Get request. Reads Instances data related to a specific project
func (a *App) GetProjectInstance(w http.ResponseWriter, r *http.Request) {

}

// Get request. Reads information about user (group)
func (a *App) GetUser(w http.ResponseWriter, r *http.Request) {

}

// Put request. Updates project with JSON payload specified
func (a *App) UpdateProject(w http.ResponseWriter, r *http.Request) {

}

// Put request. Updates a Instances associated with the project with
// the JSON payload
func (a *App) UpdateProjectInstance(w http.ResponseWriter, r *http.Request) {

}

// Put request. Updates information about user.
func (a *App) UpdateUser(w http.ResponseWriter, r *http.Request) {

}
