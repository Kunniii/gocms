package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserID   uint
	Tags     []Tag `gorm:"many2many:posts_tags;"`
	Comments []Comment
	Title    string
	Like     uint
	Body     string
}
