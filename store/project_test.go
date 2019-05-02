package store

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestProject(t *testing.T) {
	projects, err := S.ProjectList()
	assert.Nil(t, err)
	assert.Len(t, projects, 1)

	user, err := S.UserGet(username, false)
	assert.Nil(t, err)
	err = user.Validate()
	assert.Nil(t, err)

	proj, err := S.ProjectCreate("https://github.com/kantopark/warden.git", "test_project_2", "description", *user)
	assert.Nil(t, err)
	assert.NotNil(t, proj)

	projects, err = S.ProjectList()
	assert.Nil(t, err)
	assert.Len(t, projects, 2)

	projects, err = S.ProjectListByUser(username)
	assert.Nil(t, err)
	assert.Len(t, projects, 2)

	projects, err = S.ProjectListByUser("false_name")
	assert.Nil(t, err)
	assert.Len(t, projects, 0)

	_, err = S.ProjectGetById(proj.ID)
	assert.Nil(t, err)

	_, err = S.ProjectGetById(proj.ID + 10000)
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())

	_, err = S.ProjectGetByName("test_project_2")
	assert.Nil(t, err)

	_, err = S.ProjectGetByName("test_project_3")
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())

	proj.Description = "new description"
	proj, err = S.ProjectUpdate(proj)
	assert.Nil(t, err)
	assert.Equal(t, proj.Description, "new description")

	proj.ID = 0
	_, err = S.ProjectUpdate(proj)
	assert.EqualError(t, err, "id of project to update must be specified")

	err = S.ProjectDelete(proj.Name, *user)
	assert.Nil(t, err)

	projects, err = S.ProjectList()
	assert.Nil(t, err)
	assert.Len(t, projects, 1)
}
