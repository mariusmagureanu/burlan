package dao

import (
	"github.com/mariusmagureanu/burlan/src/pkg/entities"
)

type GroupDao interface {
	Insert(*entities.Group) error
	Update(*entities.Group) error
	Delete(*entities.Group) error
	GetByID(*entities.Group, uint) error
	AddUser(uint, *entities.Group) error
	RemoveUser(uint, *entities.Group) error
}

type groupDao struct {
	base
}

func (gd groupDao) Insert(g *entities.Group) error {
	return gd.db.Create(g).Error
}

func (gd groupDao) Update(g *entities.Group) error {
	return gd.db.Where("ID=?", g.ID).Save(g).Error
}

func (gd groupDao) Delete(g *entities.Group) error {
	return gd.db.Delete(g, g.ID).Error
}

func (gd groupDao) GetByID(g *entities.Group, groupID uint) error {
	return gd.db.Preload("Users").First(g, groupID).Error
}

func (gd groupDao) AddUser(groupID uint, user *entities.Group) error {
	var group entities.Group

	err := gd.GetByID(&group, groupID)

	if err != nil {
		return err
	}

	return gd.db.Model(&group).Association("Users").Append(user)
}

func (gd groupDao) RemoveUser(groupID uint, user *entities.Group) error {
	var group entities.Group

	err := gd.GetByID(&group, groupID)

	if err != nil {
		return err
	}

	return gd.db.Model(&group).Association("Users").Delete(user)
}
