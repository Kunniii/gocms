package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	PostID uint
	Name   string `gorm:"unique"`
}
