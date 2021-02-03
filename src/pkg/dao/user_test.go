package dao

import (
	"github.com/mariusmagureanu/burlan/src/pkg/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserDao_Insert(t *testing.T) {
	dao, err := tearUp()
	assert.Nil(t, err)
	defer tearDown(dao)

	user := entities.User{Name: "foo", Alias: "foo_alias", Email: "foo@email.com"}
	err = dao.Users().Insert(&user)

	assert.Nil(t, err)
	assert.Equal(t, uint(1), user.ID)
}

func TestUserDao_Delete(t *testing.T) {
	dao, err := tearUp()
	assert.Nil(t, err)
	defer tearDown(dao)

	user := &entities.User{Name: "foo", Alias: "foo_alias", Email: "foo@email.com"}
	err = dao.Users().Insert(user)
	assert.Nil(t, err)

	err = dao.Users().Delete(user)
	assert.Nil(t, err)

	foundUser := &entities.User{}
	err = dao.Users().GetByID(foundUser, user.ID)

	assert.NotNil(t, err)
}

func TestUserDao_AddFriend(t *testing.T) {
	dao, err := tearUp()
	assert.Nil(t, err)
	defer tearDown(dao)

	user := &entities.User{Name: "foo", Alias: "foo_alias", Email: "foo@email.com"}
	guest := &entities.User{Name: "baz", Alias: "baz_alias", Email: "baz@email.com"}

	err = dao.Users().Insert(user)
	assert.Nil(t, err)

	err = dao.Users().Insert(guest)
	assert.Nil(t, err)

	err = dao.Users().AddFriend(user.ID, guest.ID)
	assert.Nil(t, err)

	var foundUser entities.User
	err = dao.Users().GetByID(&foundUser, user.ID)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(foundUser.Friends))
	friend := foundUser.Friends[0]

	assert.Equal(t, guest.ID, friend.ID)
	assert.Equal(t, guest.Name, friend.Name)
	assert.Equal(t, guest.Alias, friend.Alias)
	assert.Equal(t, guest.Email, friend.Email)
}

func TestUserDao_Update(t *testing.T) {
	dao, err := tearUp()
	assert.Nil(t, err)
	defer tearDown(dao)

	user := &entities.User{Name: "foo", Alias: "foo_alias", Email: "foo@email.com"}
	err = dao.Users().Insert(user)
	assert.Nil(t, err)

	foundUser := &entities.User{}
	err = dao.Users().GetByID(foundUser, user.ID)

	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Name, foundUser.Name)
	assert.Equal(t, user.Alias, foundUser.Alias)

	foundUser.Alias = "updated alias"
	err = dao.Users().Update(foundUser)

	var updatedUser entities.User
	err = dao.Users().GetByID(&updatedUser, user.ID)

	assert.Equal(t, updatedUser.Alias, foundUser.Alias)
}

func TestUserDao_AddNonExistentFriend(t *testing.T) {
	dao, err := tearUp()
	assert.Nil(t, err)
	defer tearDown(dao)

	user := &entities.User{Name: "foo", Alias: "foo_alias", Email: "foo@email.com"}
	err = dao.Users().Insert(user)
	assert.Nil(t, err)

	err = dao.Users().AddFriend(user.ID, uint(100))
	assert.NotNil(t, err)
}

func TestUserDao_RemoveNonExistentFriend(t *testing.T) {
	dao, err := tearUp()
	assert.Nil(t, err)
	defer tearDown(dao)

	user := &entities.User{Name: "foo", Alias: "foo_alias", Email: "foo@email.com"}
	err = dao.Users().Insert(user)
	assert.Nil(t, err)

	err = dao.Users().RemoveFriend(user.ID, uint(100))
	assert.NotNil(t, err)
}

func TestUserDao_RemoveFriend(t *testing.T) {
	dao, err := tearUp()
	assert.Nil(t, err)
	defer tearDown(dao)

	user := &entities.User{Name: "foo", Alias: "foo_alias", Email: "foo@email.com"}
	guest := &entities.User{Name: "baz", Alias: "baz_alias", Email: "baz@email.com"}

	err = dao.Users().Insert(user)
	assert.Nil(t, err)

	err = dao.Users().Insert(guest)
	assert.Nil(t, err)

	err = dao.Users().AddFriend(user.ID, guest.ID)
	assert.Nil(t, err)

	var foundUser entities.User
	err = dao.Users().GetByID(&foundUser, user.ID)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(foundUser.Friends))

	err = dao.Users().RemoveFriend(foundUser.ID, guest.ID)
	assert.Nil(t, err)

	var f entities.User
	err = dao.Users().GetByID(&f, user.ID)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(f.Friends))
}

func TestUserDao_GetByID(t *testing.T) {
	dao, err := tearUp()
	assert.Nil(t, err)
	defer tearDown(dao)

	user := &entities.User{Name: "foo", Alias: "foo_alias", Email: "foo@email.com"}
	err = dao.Users().Insert(user)

	foundUser := &entities.User{}
	err = dao.Users().GetByID(foundUser, 1)
	assert.Nil(t, err)
}
