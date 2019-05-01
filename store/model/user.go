package model

import (
	"github.com/pkg/errors"

	"warden/utils"
)

type User struct {
	Id       uint      `gorm:"primary_key"`
	Username string    `gorm:"type:varchar(50);unique;index"`
	Password string    `gorm:"type:varchar(255)"`
	Projects []Project `gorm:"many2many:user_project;"`
	Type     string    `gorm:"type:varchar(10);default:'basic'"`
}

func (u *User) Validate() error {
	u.Username = utils.StrLowerTrim(u.Username)
	if u.Username == "" {
		return errors.New("Username cannot be empty")
	}

	if len(u.Password) == 0 {
		return errors.New("Password length must be > 0")
	}

	u.Type = utils.StrLowerTrim(u.Type)
	if utils.StrIsIn(u.Type, []string{"basic", "admin"}) {
		return errors.Errorf("Unknown user type: '%s'", u.Type)
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
