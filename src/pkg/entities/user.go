package entities

import (
	"gorm.io/gorm"
)

// User defines a user which will send messages.
type User struct {
	gorm.Model
	Name  string
	Alias string
	Email string
}
