package dao

import (
	"github.com/mariusmagureanu/burlan/src/pkg/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGroupDao_Insert(t *testing.T) {
	dao, err := tearUp()
	assert.Nil(t, err)
	defer tearDown(dao)

	group := entities.Group{Name: "foo"}
	err = dao.Groups().Insert(&group)

	assert.Nil(t, err)
	assert.Equal(t, uint(1), group.ID)
}

func TestGroupDao_GetByID(t *testing.T) {
	dao, err := tearUp()
	assert.Nil(t, err)
	defer tearDown(dao)

	group := &entities.Group{Name: "foo"}
	err = dao.Groups().Insert(group)

	foundGroup := &entities.Group{}
	err = dao.Groups().GetByID(foundGroup, 1)
	assert.Nil(t, err)
}
