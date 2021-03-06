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
	Description string     `gorm:"type:varchar(512)"`
	UniqueName  string     `gorm:"column:unique_name;type:varchar(100);unique;not null;index"`
	Instances   []Instance `gorm:"foreignkey:ProjectID"` // must at least have one Instance. To run the latest
	Owners      []User     `gorm:"many2many:user_project"`
}

func (p *Project) HasOwner(username string) bool {
	for _, u := range p.Owners {
		if strings.EqualFold(u.UniqueName, p.GetUniqueName(username)) {
			return true
		}
	}
	return false
}

func (p *Project) GetUniqueName(name string) string {
	return utils.StrLowerTrim(name)
}

func (p *Project) Validate() error {
	p.GitURL = strings.TrimSpace(p.GitURL)
	if matched, _ := regexp.MatchString(`^(?i)https?://\S+$`, p.GitURL); !matched {
		return errors.Errorf("GitURL: '%s' is not a valid url", p.GitURL)
	}

	p.Name = strings.TrimSpace(p.Name)
	if len(p.Name) <= 3 {
		return errors.New("Project name must be 4 characters or longer")
	}

	p.UniqueName = p.GetUniqueName(p.Name)
	return nil
}
