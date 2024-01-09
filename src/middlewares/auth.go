package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/Kunniii/gocms/internal"
	itypes "github.com/Kunniii/gocms/internal/types"
	"github.com/gin-gonic/gin"
)

func CheckAuth(context *gin.Context) {
	authorization := context.GetHeader("Authorization")

	if authorization == "" || !strings.HasPrefix(authorization, "Bearer") {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"OK":      false,
			"message": "Unauthorized",
		})
		return
	}

	if _, ok, err := internal.VerifyToken(authorization); !ok {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"OK":      false,
			"message": "Unauthorized",
		})
		return
	} else {
		context.Set("auth-token", authorization)
		context.Next()
	}

}

func RequireAdmin(context *gin.Context) {
	authorization := context.GetHeader("Authorization")

	if authorization == "" || !strings.HasPrefix(authorization, "Bearer") {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"OK":      false,
			"message": "Unauthorized",
		})
		return
	}

	if _, ok, err := internal.VerifyToken(authorization); !ok {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"OK":      false,
			"message": "Unauthorized",
		})
		return
	}

	claims := internal.GetClaims(authorization)
	roleID := uint(claims["RoleID"].(float64))

	if itypes.Roles[roleID].Name != "admin" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"OK":      false,
			"message": "Unauthorized",
		})
		return
	}

	context.Next()

}
