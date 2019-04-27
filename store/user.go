package store

import (
	"github.com/pkg/errors"

	"warden/store/model"
)

// Creates a user with the given username. Returns an error if creation fails
func (s *Store) UserCreate(name string) (*model.User, error) {
	user := &model.User{Username: name}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, errors.Wrapf(err, "could not create user with name: %s", name)
	}
	return user, nil
}

// Deletes user with specified name. Returns an error if removal fails
func (s *Store) UserDelete(name string) error {
	user, err := s.UserGet(name)
	if err != nil {
		return err
	}
	if err := s.db.Delete(user).Error; err != nil {
		return errors.Wrapf(err, "could not delete user with name: %s", name)
	}
	return nil
}

// Gets a user with specified name. Returns an error if user cannot be found
func (s *Store) UserGet(name string) (user *model.User, err error) {
	if err = s.db.Preload(_PROJECT).First(&user, "Username = ?", name).Error; err != nil {
		return nil, errors.Wrapf(err, "could not find user with name: %s", name)
	}
	return
}

// Lists all users. If withProjects is true, the projects that the users are involved in
// will be returned in the request as well. Returns an error if query fails
func (s *Store) UserList(withProjects bool) (users []*model.User, err error) {
	if withProjects {
		if err = s.db.Preload(_PROJECT).Find(&users).Error; err != nil {
			return nil, errors.Wrap(err, "could not load users")
		}
	} else {
		if err = s.db.Find(&users).Error; err != nil {
			return nil, errors.Wrapf(err, "could not load users")
		}
	}
	return
}

// Updates user with specified old name to new name. Returns an error if update fails
func (s *Store) UserUpdate(oldName, newName string) (*model.User, error) {
	user, err := s.UserGet(oldName)
	if err != nil {
		return nil, err
	}
	if _, err := s.UserGet(newName); err == nil {
		return nil, errors.Errorf("user with name '%s' already exists", newName)
	}
	user.Username = newName

	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := s.db.Save(user).Error; err != nil {
		return nil, errors.Wrapf(err, "could not update user '%s'", oldName)
	}
	return user, nil
}
