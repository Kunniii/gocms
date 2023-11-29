package itypes

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	jwt.Claims
	UserID   uint
	RoleID   uint
	Email    string
	UserName string
}
