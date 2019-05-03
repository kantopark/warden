package store

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"warden/store/model"
)

// Creates a new project. Returns an error if creation fails
func (s *Store) ProjectCreate(gitUrl, name, description string, user model.User) (*model.Project, error) {
	project := &model.Project{
		GitURL:      gitUrl,
		Name:        name,
		Description: description,
		Owners:      []model.User{user},
	}

	if err := project.Validate(); err != nil {
		return nil, err
	}

	if err := s.db.Create(project).Error; err != nil {
		return nil, errors.Wrapf(err, "error creating project")
	}

	return project, nil
}

// Deletes a project from the database. Returns an error if removal fails
func (s *Store) ProjectDelete(name string) error {
	project, err := s.ProjectGetByName(name)
	if err != nil {
		return err
	}

	if err := s.db.Delete(project).Error; err != nil {
		return errors.Wrapf(err, "error removing project")
	}
	return nil
}

// Searches for a project by it's ID. Returns an error if query fails
func (s *Store) ProjectGetById(id uint) (*model.Project, error) {
	var project model.Project
	if err := s.db.Preload(_OWNERS).Preload(_INSTANCES).First(&project, id).Error; err == gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil {
		return nil, errors.Wrapf(err, "could not find project with id: %d", id)
	}
	return &project, nil
}

// Searches for a project by it's name. Returns an error if query fails
func (s *Store) ProjectGetByName(name string) (*model.Project, error) {
	var project model.Project
	if err := s.db.Preload(_OWNERS).Preload(_INSTANCES).First(&project, "unique_name = ?", project.GetUniqueName(name)).Error; err == gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil {
		return nil, errors.Wrapf(err, "could not find project with name: %s", name)
	}
	return &project, nil
}

// Lists all projects. This is normally used by admins since it'll list all projects
func (s *Store) ProjectList() (projects []model.Project, err error) {
	if err := s.db.Preload(_OWNERS).Preload(_INSTANCES).Find(&projects).Error; err != nil {
		return nil, errors.Wrap(err, "could not list projects")
	}
	return
}

func (s *Store) ProjectListByUser(username string) (projects []model.Project, err error) {
	proj, err := s.ProjectList()
	if err != nil {
		return nil, err
	}
	for _, p := range proj {
		if p.HasOwner(username) {
			projects = append(projects, p)
		}
	}

	for _, p := range projects {
		log.Printf("%+v\n", p)
	}
	return
}

// Updates the existing project with the new project. The existing project is identified by the ID
// which must be specified in the newProj argument. Returns an error if update fails
func (s *Store) ProjectUpdate(newProj *model.Project) (*model.Project, error) {
	if newProj.ID == 0 {
		return nil, errors.New("id of project to update must be specified")
	}

	project, err := s.ProjectGetById(newProj.ID)
	if err != nil {
		return nil, err
	}

	project.Name = newProj.Name
	project.UniqueName = newProj.GetUniqueName(project.Name)
	project.Description = newProj.Description
	project.GitURL = newProj.GitURL
	if newProj.Owners != nil && len(newProj.Owners) > 0 {
		project.Owners = newProj.Owners
	}

	if err := project.Validate(); err != nil {
		return nil, err
	}

	if err := s.db.Save(project).Error; err != nil {
		return nil, errors.Wrap(err, "could not update project")
	}

	return project, nil
}
