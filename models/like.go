package models

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	ID     uint
	PostID uint
	UserID uint
}
