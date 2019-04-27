package model

import (
	"strings"

	"github.com/pkg/errors"
)

type User struct {
	Id       uint      `gorm:"primary_key"`
	Username string    `gorm:"type:varchar(50);unique;index"`
	Projects []Project `gorm:"many2many:user_project;"`
}

func (u *User) Validate() error {
	u.Username = strings.TrimSpace(u.Username)
	if u.Username == "" {
		return errors.New("Username cannot be empty")
	}
	return nil
}
