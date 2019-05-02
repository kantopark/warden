package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore_UserGet(t *testing.T) {
	user, err := S.UserGet(username, false)
	assert.Nil(t, err)
	assert.Len(t, user.Projects, 1)
	assert.False(t, user.IsAdmin())
}

func TestStore_User(t *testing.T) {
	users, err := S.UserList(false)
	assert.Nil(t, err)
	assert.Len(t, users, 1)

	newUser, err := S.UserCreate(UserBody{
		Email:    "admin@admin.com",
		Username: "admin",
		Password: "admin",
	})
	assert.Nil(t, err)
	assert.NotNil(t, newUser)

	newUser, err = S.UserGet("admin", false)
	assert.Nil(t, err)
	assert.NotNil(t, newUser)

	_, err = S.UserLogin(UserBody{
		Email:    "admin@admin.com",
		Username: "admin",
		Password: "bad_password",
	}, true)
	assert.NotNil(t, err)

	_, err = S.UserLogin(UserBody{
		Email:    "admin@admin.com",
		Username: "admin",
		Password: "admin",
	}, true)
	assert.Nil(t, err)

	newUser, err = S.UserUpdate(UserBody{
		Email:       "admin2@admin.com",
		Username:    "admin",
		Password:    "admin",
		NewPassword: "new_password",
	})
	assert.Nil(t, err)
	assert.Equal(t, newUser.Email, "admin2@admin.com")
	assert.Equal(t, newUser.Password, "new_password")

	users, err = S.UserList(false)
	assert.Nil(t, err)
	assert.Len(t, users, 2)

	err = S.UserDelete("admin")
	assert.Nil(t, err)

	users, err = S.UserList(false)
	assert.Nil(t, err)
	assert.Len(t, users, 1)
}
