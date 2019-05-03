package application

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"warden/store/model"
	"warden/utils"
)

// Post request. Appends a new Instances to Project.
func (a *App) CreateProjectInstance(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)

	var i model.Instance
	if err := parseJson(r.Body, &i); err != nil {
		internalServerError(w, errors.Wrap(err, "error parsing JSON"))
		return
	}

	proj, err := a.db.ProjectGetById(i.ProjectID)
	if err != nil {
		internalServerError(w, errors.Wrapf(err, "error getting project with ID = %d", i.ProjectID))
		return
	}
	if !proj.HasOwner(u.Username) {
		forbidden(w, errors.New("you're not authorized to make changes to this project"))
		return
	}

	inst, err := a.db.InstanceCreate(i.CommitHash, i.Alias, proj.Name)
	if err != nil {
		internalServerError(w, errors.Wrap(err, "error creating instance"))
	}
	jsonify(w, inst)
}

// Delete request. Removes an instance from the project. Project must however be left with
// at least the latest instance.
func (a *App) DeleteProjectInstance(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	projectName := chi.URLParam(r, "name")
	id, err := strconv.Atoi(utils.StrLowerTrim(chi.URLParam(r, "id")))
	if err != nil {
		badRequest(w, errors.New("unable to parse id field as an integer"))
		return
	} else if id <= 0 {
		badRequest(w, errors.New("instance id must be >= 0"))
		return
	}

	proj, err := a.db.ProjectGetByName(projectName)
	if err != nil {
		internalServerError(w, errors.Wrapf(err, "error getting project with name = %s", projectName))
		return
	}
	if !proj.HasOwner(u.Username) {
		forbidden(w, errors.New("you're not authorized to make changes to this project"))
		return
	}

	for _, i := range proj.Instances {
		if i.ID == uint(id) {
			if err := a.db.InstanceDelete(i.ProjectID, i.CommitHash); err != nil {
				internalServerError(w, errors.Wrap(err, "error removing instance"))
				return
			}
			ok(w)
			return
		}

	}

	badRequest(w, errors.Errorf("could not find instance with project name '%s' and id '%d'", projectName, id))
}

// Put request. Updates a Instances associated with the project with
// the JSON payload
func (a *App) UpdateProjectInstance(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)

	var i model.Instance
	if err := parseJson(r.Body, &i); err != nil {
		internalServerError(w, errors.Wrap(err, "error parsing JSON"))
		return
	}
	if i.ID <= 0 {
		badRequest(w, errors.New("instnace id must be >= 0"))
		return
	}

	proj, err := a.db.ProjectGetById(i.ProjectID)
	if err != nil {
		internalServerError(w, errors.Wrapf(err, "error getting project with id = %d", i.ProjectID))
		return
	}
	if !proj.HasOwner(u.Username) {
		forbidden(w, errors.New("you're not authorized to make changes to this project"))
		return
	}

	for _, inst := range proj.Instances {
		if inst.ID == i.ID {
			updated_inst, err := a.db.InstanceUpdate(&i)
			if err != nil {
				internalServerError(w, errors.Wrap(err, "could not update instance"))
				return
			}
			jsonify(w, updated_inst)
			return
		}
	}
	badRequest(w, errors.Errorf("could not find instance with project id '%d' and instance id '%d'", i.ProjectID, i.ID))
}
