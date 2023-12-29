package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	ID     uint
	PostID uint
	UserID uint
	Like   uint
	Body   string
}
