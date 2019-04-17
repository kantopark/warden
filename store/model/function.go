package model

import "github.com/jinzhu/gorm"

// The Function object. This model stores information such as the name,
// description and git url. The git url specifies where to download the
// function code from. The specific runtime information such as the
// specific commit to run are given in the RunInfo
type Function struct {
	*gorm.Model
	GitURL      string    `gorm:"column:git_url"`
	Name        string    `gorm:"type:varchar(100)"`
	Description string    `gorm:"type:varchar(512)"`
	RunInfo     []RunInfo // must at least have one RunInfo. To run the latest
}
