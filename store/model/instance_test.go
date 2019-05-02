package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstance(t *testing.T) {
	inst := &Instance{
		ID:         1,
		Alias:      "",
		CommitHash: "129841284912812849081204",
		ProjectID:  1,
	}

	err := inst.Validate()
	assert.Nil(t, err)
	assert.Equal(t, inst.Alias, "latest")

	inst.Alias = "  Dev  "
	err = inst.Validate()
	assert.Nil(t, err)
	assert.Equal(t, inst.Alias, "dev")

	inst.ProjectID = 0
	err = inst.Validate()
	assert.EqualError(t, err, "runtime instance must be linked to a project instance via a project id key")

	inst.CommitHash = ""
	err = inst.Validate()
	assert.EqualError(t, err, "commit hash for runtime instance cannot be empty")
}
