package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	PostID uint
	UserID uint
	Like   uint
	Body   string
}
