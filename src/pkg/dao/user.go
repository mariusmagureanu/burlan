package dao

import (
	"github.com/mariusmagureanu/burlan/src/pkg/entities"
)

type UserDao interface {
	Insert(*entities.User) error
	Update(*entities.User) error
	Delete(*entities.User) error
	GetByID(*entities.User, uint) error
	AddFriend(uint, uint) error
	RemoveFriend(uint, uint) error
}

type userDao struct {
	base
}

func (ud userDao) Insert(u *entities.User) error {
	return ud.db.Create(u).Error
}

func (ud userDao) Update(u *entities.User) error {
	return ud.db.Where("ID=?", u.ID).Save(u).Error
}

func (ud userDao) Delete(u *entities.User) error {
	return ud.db.Delete(u, u.ID).Error
}

func (ud userDao) GetByID(u *entities.User, userID uint) error {
	return ud.db.Preload("Friends").First(u, userID).Error
}

func (ud userDao) AddFriend(userId uint, friendId uint) error {
	var (
		user   entities.User
		friend entities.User
		err    error
	)

	err = ud.GetByID(&user, userId)

	if err != nil {
		return err
	}

	err = ud.db.First(&friend, friendId).Error

	if err != nil {
		return err
	}

	err = ud.db.Model(&user).Association("Friends").Append(&friend)

	return err
}

func (ud userDao) RemoveFriend(userId uint, friendId uint) error {
	var (
		user   entities.User
		friend entities.User
		err    error
	)

	err = ud.GetByID(&user, userId)

	if err != nil {
		return err
	}

	err = ud.db.First(&friend, friendId).Error

	if err != nil {
		return err
	}

	err = ud.db.Model(&user).Association("Friends").Delete(&friend)

	return err
}
