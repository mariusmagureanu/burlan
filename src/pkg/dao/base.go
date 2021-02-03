package dao

import "gorm.io/gorm"

type base struct {
	db *gorm.DB
}
