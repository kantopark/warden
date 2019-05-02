package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProject(t *testing.T) {
	user := User{
		ID:         1,
		Email:      "daniel.bok@outlook.com",
		Username:   "daniel",
		UniqueName: "daniel",
		Password:   "password",
		Type:       "basic",
	}

	project := &Project{
		ID:          1,
		GitURL:      "https://github.com/yi-jiayu/bus-eta-bot.git",
		Name:        "BusEta",
		Description: "Jiayu's bus ETA bot",
		Instances:   nil,
		Owners:      []User{user},
	}

	err := project.Validate()
	assert.Nil(t, err)
	assert.Equal(t, project.UniqueName, "buseta")

	assert.True(t, project.HasOwner(user))

	user.Username = "bad_username"
	assert.False(t, project.HasOwner(user))

	project.GitURL = "http://github.com/yi-jiayu/bus-eta-bot.git"
	err = project.Validate()
	assert.Nil(t, err)

	project.Name = "Bus"
	err = project.Validate()
	assert.EqualError(t, err, "Project name must be 4 characters or longer")

	project.GitURL = "github.com/yi-jiayu/bus-eta-bot.git"
	err = project.Validate()
	assert.EqualErrorf(t, err, "GitURL: 'github.com/yi-jiayu/bus-eta-bot.git' is not a valid url", "GitURL: '%s' is not a valid url", project.GitURL)
}
