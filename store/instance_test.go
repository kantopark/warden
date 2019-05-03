package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstance(t *testing.T) {
	inst, err := S.InstanceGetByHash(1, "95bfc3515452bfafeb2e04f948ac26d1e2a871c8")
	assert.Nil(t, err)

	proj, err := S.ProjectGetById(inst.ProjectID)
	assert.Nil(t, err)

	inst, err = S.InstanceCreate("0c0aafa7ec1250be737d0d39f6de36854baa0f8b", "", proj.Name)
	assert.Nil(t, err)

	inst.Alias = "test-2"
	inst, err = S.InstanceUpdate(inst)
	assert.Nil(t, err)

	err = S.InstanceDelete(inst.ProjectID, inst.CommitHash)
	assert.Nil(t, err)
}
