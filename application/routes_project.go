package application

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

type projectBody struct {
	Description string `json:"description"`
	GitURL      string `json:"git_url"`
	Name        string `json:"name"`
}

// Post request. Creates a new project in the system. JSON payload
// should specify creation parameters
func (a *App) CreateProject(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)

	var p projectBody
	if err := parseJson(r.Body, &p); err != nil {
		internalServerError(w, errors.Wrap(err, "error parsing JSON"))
		return
	}

	user, err := a.db.UserGet(u.Username, true)
	if err != nil {
		internalServerError(w, errors.Wrap(err, "error getting user info"))
		return
	}

	project, err := a.db.ProjectCreate(p.GitURL, p.Name, p.Description, *user)
	if err != nil {
		internalServerError(w, errors.Wrap(err, "error creating project"))
	}
	jsonify(w, project)
}

// Delete request. Removes project and all Instances associated with it
func (a *App) DeleteProject(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	proj, err := a.db.ProjectGetByName(chi.URLParam(r, "name"))
	if err != nil {
		internalServerError(w, errors.Wrap(err, "error getting project"))
		return
	}

	if !proj.HasOwner(u.Username) {
		forbidden(w, errors.Wrap(err, "you're not authorized to delete this project"))
		return
	}

	if err := a.db.ProjectDelete(proj.Name); err != nil {
		internalServerError(w, errors.Wrap(err, "error deleting project"))
		return
	}

	ok(w)
}

// Get request. Reads information related to a specific project
func (a *App) GetProject(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)

	proj, err := a.db.ProjectGetByName(chi.URLParam(r, "name"))
	if err != nil {
		internalServerError(w, errors.Wrap(err, "error getting project"))
		return
	}
	if !proj.HasOwner(u.Username) {
		forbidden(w, errors.Wrap(err, "you're not authorized to view this project"))
		return
	}

	jsonify(w, proj)
}

// Get request. List all project
func (a *App) ListProjects(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	projects, err := a.db.ProjectListByUser(u.Username)
	if err != nil {
		internalServerError(w, err)
		return
	}
	jsonify(w, projects)
}

// Put request. Updates project with JSON payload specified
func (a *App) UpdateProject(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)

	var p projectBody
	if err := parseJson(r.Body, &p); err != nil {
		internalServerError(w, errors.Wrap(err, "error parsing JSON"))
		return
	}
	proj, err := a.db.ProjectGetByName(p.Name)
	if err != nil {
		internalServerError(w, errors.Wrapf(err, "error getting project with name '%s'", p.Name))
		return
	}
	if !proj.HasOwner(u.Username) {
		forbidden(w, errors.Wrap(err, "you're not authorized to update this project"))
		return
	}

	proj.GitURL = p.GitURL
	proj.Description = p.Description
	proj.Name = p.Name
	if err := proj.Validate(); err != nil {
		badRequest(w, errors.Wrap(err, "invalid update parameters for project"))
		return
	}

	proj, err = a.db.ProjectUpdate(proj)
	if err != nil {
		internalServerError(w, errors.Wrap(err, "could not update project"))
		return
	}

	jsonify(w, proj)
}
