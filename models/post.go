package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserID   uint
	Tags     []*Tag     `gorm:"many2many:posts_tags;"`
	Comments []*Comment `gorm:"many2many:posts_comments;"`
	Title    string
	Like     uint
	Body     string
}
