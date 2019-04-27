package store

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"warden/store/model"
)

// Creates a runnable instance for the project
func (s *Store) InstanceCreate(commit, alias string, projectName string) (*model.Instance, error) {
	project, err := s.ProjectGetByName(projectName)
	if err != nil {
		return nil, err
	}
	instance := &model.Instance{
		Alias:      alias,
		CommitHash: commit,
		ProjectID:  project.ID,
	}

	if err = instance.Validate(); err != nil {
		return nil, err
	}

	return instance, nil
}

// Gets a running instance of the project by the project ID and the commit hash
func (s *Store) InstanceGetByHash(projectID uint, commitHash string) (inst *model.Instance, err error) {
	if err = s.db.First(&inst, "ProjectID = ? AND CommitHash = ?", projectID, commitHash).Error; err == gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil {
		return nil, errors.Wrapf(err, "error getting instance with project id '%d' and commit hash '%s'", projectID, commitHash)
	}
	return
}

// Gets a running instance of the project by the instance ID
func (s *Store) InstanceGetById(id uint) (inst *model.Instance, err error) {
	if err := s.db.First(&inst, id).Error; err == gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil {
		return nil, errors.Wrapf(err, "error getting instance with id '%d' and alias '%s'", id)
	}
	return
}

// Deletes a running instance of the project
func (s *Store) InstanceDelete(projectID uint, commitHash string) error {
	project, err := s.InstanceGetByHash(projectID, commitHash)
	if err != nil {
		return err
	}
	if err := s.db.Delete(project).Error; err != nil {
		return errors.Wrapf(err, "error removing instance with project '%d' and commit hash '%s'", projectID, commitHash)
	}
	return nil
}

// Updates a running instance of the project by the instance ID. Since it is an update, it assumes
// that the user already has the ID of the instance. Thus we search for existing instance by the
// instance ID
func (s *Store) InstanceUpdateByAlias(newInstance *model.Instance) (*model.Instance, error) {
	if err := newInstance.Validate(); err != nil {
		return nil, err
	}
	inst, err := s.InstanceGetById(newInstance.Id)
	if err != nil {
		return nil, err
	}
	inst.CommitHash = newInstance.CommitHash
	inst.Alias = newInstance.Alias

	if err := s.db.Save(inst).Error; err != nil {
		return nil, errors.Wrapf(err, "could not update instance: %+v", inst)
	}

	return inst, nil
}
