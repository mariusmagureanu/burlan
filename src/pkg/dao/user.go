package dao

import (
	"github.com/mariusmagureanu/burlan/src/pkg/entities"

	"github.com/google/uuid"
)

type UserDao interface {
	Insert(*entities.User) error
	Update(*entities.User) error
	Delete(*entities.User) error
	GetByID(*entities.User, uint) error
	GetByName(*entities.User, string) error
	GetByUID(*entities.User, string) error
	AddFriend(uint, uint) error
	RemoveFriend(uint, uint) error
	GetAll(*[]entities.User) error
}

type userDao struct {
	base
}

func (ud userDao) GetAll(users *[]entities.User) error {
	return ud.db.Find(users).Error
}

func (ud userDao) Insert(u *entities.User) error {
	u.UID = uuid.NewString()
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

func (ud userDao) GetByName(u *entities.User, userName string) error {
	return ud.db.Preload("Friends").Where("name=?", userName).First(u).Error
}

func (ud userDao) GetByUID(u *entities.User, uid string) error {
	return ud.db.Preload("Friends").Where("uid=?", uid).First(u).Error
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
