package store

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	"warden/store/model"
)

var hashCost int

func init() {
	hashCost = viper.GetInt("auth.hash_cost")
	if hashCost < 10 {
		hashCost = 10
	}
}

// The user body payload sent from the front end
type UserBody struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password,omitempty"`
}

// Creates a user with the given username. Returns an error if creation fails
func (s *Store) UserCreate(u *UserBody) (*model.User, error) {
	pw, err := HashPassword(u.Password)
	if err != nil {
		return nil, err
	}
	user := &model.User{Username: u.Username, Password: string(pw)}
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, errors.Wrapf(err, "could not create user with name: %s", user.Username)
	}
	return user, nil
}

// Deletes user with specified name. Returns an error if removal fails
func (s *Store) UserDelete(u *UserBody) error {
	user, err := s.UserGet(u.Username, false)
	if err != nil {
		return err
	}

	if err := ComparePasswords(u.Password, user.Password); err != nil {
		return err
	}

	if err := s.db.Delete(user).Error; err != nil {
		return errors.Wrapf(err, "could not delete user with name: %s", user.Username)
	}
	return nil
}

// Gets a user with specified name. Returns an error if user cannot be found
func (s *Store) UserGet(name string, noPassword bool) (user *model.User, err error) {
	if err = s.db.Preload(_PROJECT).First(&user, "Username = ?", name).Error; err != nil {
		return nil, errors.Wrapf(err, "could not find user with name: %s", name)
	}
	if noPassword {
		user.Password = ""
	}
	return
}

// Lists all users along with the projects that the users are involved in.
// Returns an error if query fails
func (s *Store) UserList(noPassword bool) (users []*model.User, err error) {
	if err = s.db.Preload(_PROJECT).Find(&users).Error; err != nil {
		return nil, errors.Wrap(err, "could not load users")
	}
	if noPassword {
		for _, u := range users {
			u.Password = ""
		}
	}
	return
}

// Updates user with specified old name to new name. Returns an error if update fails
func (s *Store) UserUpdatePassword(u *UserBody) (*model.User, error) {
	user, err := s.UserGet(u.Username, false)
	if err != nil {
		return nil, err
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}
	if err := ComparePasswords(u.Password, user.Password); err != nil {
		return nil, err
	}

	user.Password = u.NewPassword
	if err := user.Validate(); err != nil {
		return nil, errors.Wrap(err, "new password does not fulfil criteria")
	}

	if err := s.db.Save(user).Error; err != nil {
		return nil, errors.Wrapf(err, "could not update password for user '%s'", u.Username)
	}
	return user, nil
}

// Hashes the given string (password)
func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return nil, errors.Wrap(err, "error hashing password")
	}
	return hashedPassword, nil
}

// Compares the plain-text original string (password) against the hashed password.
// After hashing the original, returns no errors if both matches and an error otherwise.
func ComparePasswords(plain, hashed string) error {
	hashedPlain, err := HashPassword(plain)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(hashedPlain, []byte(hashed))
	if err != nil {
		return errors.Wrap(err, "passwords do not match")
	}
	return nil
}
