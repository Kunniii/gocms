package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	AuthorID uint
	Author   User
	Tags     []Tag `gorm:"many2many:post_tags;"`
	Comments []Comment
	Title    string
	Like     uint
	Body     string
}
