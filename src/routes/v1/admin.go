package v1

import (
	"net/http"

	"github.com/Kunniii/gocms/controllers"
	"github.com/Kunniii/gocms/middlewares"
	"github.com/gin-gonic/gin"
)

func UseAdminRouter(router *gin.RouterGroup) {
	router.GET("/ready", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"OK": true,
		})
	})

	router.POST("/create-user", middlewares.RequireAdmin, controllers.AddUser)

}
