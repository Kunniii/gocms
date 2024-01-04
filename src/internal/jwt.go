package internal

import (
	"fmt"
	"os"
	"time"

	itypes "github.com/Kunniii/gocms/internal/types"
	"github.com/golang-jwt/jwt/v5"
)

var JWT_KEY = []byte(os.Getenv("JWT_KEY"))

func CreateToken(claims itypes.UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserID":   claims.UserID,
		"RoleID":   claims.RoleID,
		"Email":    claims.Email,
		"UserName": claims.UserName,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(JWT_KEY)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return JWT_KEY, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func VerifyToken(tokenString string) (*jwt.Token, bool, error) {
	tokenString = tokenString[len("Bearer "):]
	token, err := ParseToken(tokenString)
	if err != nil {
		return nil, false, err
	}

	if !token.Valid {
		return nil, false, fmt.Errorf("invalid token")
	}

	return token, true, nil
}

func GetClaims(authToken string) jwt.MapClaims {
	token, _, _ := VerifyToken(authToken)
	claims := token.Claims.(jwt.MapClaims)
	return claims
}
