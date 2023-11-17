package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string
	Email    string `gorm:"unique"`
	Password string
	RoleID   uint
	Role     Role
}
