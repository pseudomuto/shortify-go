package app

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type UserSuite struct {
	suite.Suite
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (suite *UserSuite) SetupSuite() {
	Configure("../examples/sqlite3.gcfg")
}

func (suite *UserSuite) TearDownTest() {
	shortifyDb.reset()
}

func (suite *UserSuite) TestNewUser() {
	t := suite.T()
	user := NewUser("testuser")

	assert.Equal(t, "testuser", user.Name)
	assert.NotEmpty(t, user.Password)
}

func (suite *UserSuite) TestGetUsers() {
	t := suite.T()
	user := NewUser("testuser")
	err := user.Save()
	assert.Nil(t, err)

	users, err := GetUsers()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
}

func (suite *UserSuite) TestUserSave() {
	t := suite.T()
	user := NewUser("testuser")
	err := user.Save()

	assert.Nil(t, err)
	assert.NotEqual(t, 0, user.Id)
	assert.NotEmpty(t, user.PasswordHash)
	assert.WithinDuration(t, time.Now(), user.CreatedAt, 100*time.Millisecond)
}

func (suite *UserSuite) TestUserNameMustBeUnique() {
	t := suite.T()
	user := NewUser("testuser")
	err := user.Save()
	assert.Nil(t, err)

	user = NewUser("testuser")
	err = user.Save()
	assert.NotNil(t, err)
}

func (suite *UserSuite) TestIsValidUser() {
	t := suite.T()
	user := NewUser("testuser")
	err := user.Save()
	assert.Nil(t, err)

	valid := IsValidUser("testuser", user.Password)
	assert.True(t, valid)

	valid = IsValidUser("testuser", "someOtherPassword")
	assert.False(t, valid)

	valid = IsValidUser("whoami", user.Password)
	assert.False(t, valid)
}

func (suite *UserSuite) TestGetUser() {
	t := suite.T()
	user := NewUser("testuser")
	err := user.Save()
	assert.Nil(t, err)

	found, err := GetUser("testuser")
	assert.Nil(t, err)
	assert.Equal(t, user.Id, found.Id)

	found, err = GetUser("whoisthis")
	assert.NotNil(t, err)
}

func (suite *UserSuite) TestResetPassword() {
	t := suite.T()
	user := NewUser("testuser")
	err := user.Save()
	assert.Nil(t, err)

	oldPassword := user.Password
	oldPasswordHash := user.PasswordHash

	err = user.ResetPassword()
	assert.Nil(t, err)
	assert.NotEqual(t, oldPassword, user.Password)
	assert.NotEqual(t, oldPasswordHash, user.PasswordHash)
}
