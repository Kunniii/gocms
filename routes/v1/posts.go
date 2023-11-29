package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UsePostsRouter(router *gin.RouterGroup) {
	router.GET("/ready", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"OK": true,
		})
	})
}