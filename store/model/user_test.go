package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	user := &User{
		ID:       1,
		Email:    "daniel.bok@outlook.com",
		Username: "   Daniel   ",
		Password: "my_password",
		Projects: nil,
		Type:     "basic",
	}

	err := user.Validate()
	assert.Nil(t, err)
	assert.Equal(t, user.UniqueName, "daniel") // should have lowered and stripped spaces
	assert.Equal(t, user.Username, "Daniel")   // should have stripped spaces

	assert.False(t, user.IsAdmin())

	user.Email = "invalid_email"
	err = user.Validate()
	assert.EqualErrorf(t, err, "invalid_email is not a valid email", "%s is not a valid email", user.Email)

	user.Type = "bad_type"
	err = user.Validate()
	assert.EqualErrorf(t, err, "Unknown user type: 'bad_type'", "Unknown user type: '%s'", user.Type)

	user.MaskPassword(true)
	assert.Equal(t, user.Password, "")
	err = user.Validate()
	assert.EqualError(t, err, "Password length must be > 0")

	user.Username = ""
	err = user.Validate()
	assert.EqualError(t, err, "Username cannot be empty")
}
