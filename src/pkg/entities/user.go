package entities

import (
	"gorm.io/gorm"
)

// User defines a user which will send messages.
type User struct {
	gorm.Model
	Name    string `gorm:"unique"`
	UID 	string `gorm:"unique"`
	Alias   string
	Email   string `gorm:"unique"`
	Friends []User `gorm:"many2many:user_friends"`
}
