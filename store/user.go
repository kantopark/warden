package store

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	"warden/store/model"
	"warden/utils"
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
	Email       string `json:"email,omitempty"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password,omitempty"`
}

// Creates a user with the given username. Returns an error if creation fails
func (s *Store) UserCreate(u UserBody) (*model.User, error) {
	pw, err := HashPassword(u.Password)
	if err != nil {
		return nil, err
	}
	user := &model.User{Username: u.Username, Password: string(pw), Email: u.Email, Type: "basic"}
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, errors.Wrapf(err, "could not create user with name: %s", user.Username)
	}
	return user, nil
}

// Deletes user with specified name. Returns an error if removal fails
// Only admin can run this function
func (s *Store) UserDelete(name string) error {
	user, err := s.UserGet(name, false)
	if err != nil {
		return err
	}

	if err := s.db.Delete(user).Error; err != nil {
		return errors.Wrapf(err, "could not delete user with name: %s", user.Username)
	}
	return nil
}

// Gets a user with specified name. Returns an error if user cannot be found
func (s *Store) UserGet(name string, maskPassword bool) (*model.User, error) {
	var user model.User
	if err := s.db.Preload(_PROJECTS).First(&user, "unique_name = ?", utils.StrLowerTrim(name)).Error; err != nil {
		return nil, errors.Wrapf(err, "could not find user with username: %s", name)
	}
	user.MaskPassword(maskPassword)
	return &user, nil
}

// Lists all users along with the projects that the users are involved in.
// Returns an error if query fails
func (s *Store) UserList(maskPassword bool) (users []model.User, err error) {
	if err = s.db.Preload(_PROJECTS).Find(&users).Error; err != nil {
		return nil, errors.Wrap(err, "could not load users")
	}

	if maskPassword {
		for _, u := range users {
			u.MaskPassword(maskPassword)
		}
	}
	return
}

// Updates user. Returns an error if update fails
func (s *Store) UserUpdate(u UserBody) (*model.User, error) {
	user, err := s.UserGet(u.Username, false)
	if err != nil {
		return nil, err
	}

	// change password
	if u.NewPassword != "" {
		if err := ComparePasswords(u.Password, user.Password); err != nil {
			return nil, err
		}
		user.Password = u.NewPassword
	}

	// change email
	if u.Email != user.Email {
		user.Email = u.Email
	}

	// check that all is well
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := s.db.Save(user).Error; err != nil {
		return nil, errors.Wrapf(err, "could not update password for user '%s'", u.Username)
	}
	return user, nil
}

// Checks if the user credentials are valid. If credentials are invalid, returns the User object.
// Otherwise, returns an error
func (s *Store) UserLogin(u UserBody, maskPassword bool) (*model.User, error) {
	user, err := s.UserGet(u.Username, false)
	if err != nil {
		return nil, err
	}
	if err := ComparePasswords(u.Password, user.Password); err != nil {
		return nil, err
	}
	user.MaskPassword(maskPassword)
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
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)); err != nil {
		return errors.Wrap(err, "passwords do not match")
	}
	return nil
}
