package entities

import (
	"gorm.io/gorm"
)

// Group represents a logical
// collection of users.
type Group struct {
	gorm.Model
	Name  string
	Users []User `gorm:"many2many:group_users;"`
}
