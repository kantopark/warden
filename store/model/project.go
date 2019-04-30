package model

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"warden/utils"
)

// The Project object. This model stores information such as the name,
// description and git url. The git url specifies where to download the
// function code from. The specific runtime information such as the
// specific commit to run are given in the Instance
type Project struct {
	ID          uint       `gorm:"primary_key"`
	GitURL      string     `gorm:"column:git_url"`
	Name        string     `gorm:"type:varchar(100)"`
	NameUnique  string     `gorm:"type:varchar(100);unique;not null;index"`
	Description string     `gorm:"type:varchar(512)"`
	Instances   []Instance // must at least have one Instance. To run the latest
	Owners      []User     `gorm:"many2many:user_project;"`
}

func (p *Project) HasOwner(user User) bool {
	for _, u := range p.Owners {
		if strings.EqualFold(u.Username, user.Username) {
			return true
		}
	}
	return false
}

func (p *Project) GetUniqueName() string {
	return utils.StrLowerTrim(p.Name)
}

func (p *Project) Validate() error {
	p.GitURL = strings.TrimSpace(p.GitURL)
	if matched, err := regexp.MatchString(`https?://\S+`, p.GitURL); !matched || err != nil {
		return errors.Wrapf(err, "GitURL: '%s' is not a valid url", p.GitURL)
	}

	p.Name = strings.TrimSpace(p.Name)
	if len(p.Name) <= 3 {
		return errors.New("Project name must be 4 characters or longer")
	}

	p.NameUnique = p.GetUniqueName()
	return nil
}
