package model

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"warden/utils"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type User struct {
	ID         uint      `gorm:"primary_key"`
	Email      string    `gorm:"type:varchar(255);unique;index"`
	Username   string    `gorm:"type:varchar(100)"`
	UniqueName string    `gorm:"column:unique_name;type:varchar(100);unique;index"`
	Password   string    `gorm:"type:varchar(255)"`
	Projects   []Project `gorm:"many2many:user_project;"`
	Type       string    `gorm:"type:varchar(10);default:'basic'"`
}

func (u *User) Validate() error {
	u.Username = strings.TrimSpace(u.Username)
	u.UniqueName = utils.StrLowerTrim(u.Username)
	if u.Username == "" {
		return errors.New("Username cannot be empty")
	}

	if len(u.Password) == 0 {
		return errors.New("Password length must be > 0")
	}

	u.Type = utils.StrLowerTrim(u.Type)
	if !utils.StrIsIn(u.Type, []string{"basic", "admin"}) {
		return errors.Errorf("Unknown user type: '%s'", u.Type)
	}

	if !emailRegex.MatchString(u.Email) {
		return errors.Errorf("%s is not a valid email", u.Email)
	}

	return nil
}

func (u *User) IsAdmin() bool {
	return utils.StrLowerTrim(u.Type) == "admin"
}

func (u *User) MaskPassword(mask bool) {
	if mask {
		u.Password = ""
	}
}
