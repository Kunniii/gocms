package middlewares

import (
	"log"
	"net/http"

	"github.com/Kunniii/gocms/internal"
	"github.com/gin-gonic/gin"
)

func CheckAuth(context *gin.Context) {
	authorization := context.GetHeader("Authorization")
	if authorization == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"OK":      false,
			"message": "Unauthorized",
		})
	}

	if token, ok, err := internal.VerifyToken(authorization); !ok {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"OK":      false,
			"message": "Unauthorized",
		})
	} else {
		context.Set("user-data", token.Claims)
		context.Next()
	}

}
