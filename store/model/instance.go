package model

import (
	"strings"

	"github.com/pkg/errors"

	"warden/utils"
)

// The Instance contains information on the how to run an instance of
// the project. Specifically, it links the alias to the commit. Each
// instance is akin to running the specific commit hash of the function
// The default alias is the "latest" tag, thus, it would mean a request
// to the base url of the function will hit that specific runtime.
// Internally, the default empty string "" will be aliased to "latest"
// as well.
type Instance struct {
	ID         uint   `gorm:"primary_key"`
	Alias      string `gorm:"unique_index:idx_alias_function"`
	CommitHash string `gorm:"column:commit_hash;varchar(100)"`
	ProjectID  uint   `gorm:"unique_index:idx_alias_function"`
}

func (i *Instance) Validate() error {
	i.Alias = utils.StrLowerTrim(i.Alias)
	if i.Alias == "" {
		i.Alias = "latest"
	}

	i.CommitHash = strings.TrimSpace(i.CommitHash)
	if i.CommitHash == "" {
		return errors.New("commit hash for runtime instance cannot be empty")
	}

	if i.ProjectID == 0 {
		return errors.New("runtime instance must be linked to a project instance via a project id key")
	}
	return nil
}
