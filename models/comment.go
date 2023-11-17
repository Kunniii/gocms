package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Post     Post
	PostID   uint
	AuthorID uint
	Author   User
	Like     uint
	Body     string
}
