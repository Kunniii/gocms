package v1

import (
	"net/http"

	"github.com/Kunniii/gocms/controllers"
	"github.com/Kunniii/gocms/middlewares"
	"github.com/gin-gonic/gin"
)

func UseTagRouter(router *gin.RouterGroup) {
	router.GET("/ready", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"OK": true,
		})
	})

	router.POST("/", middlewares.CheckAuth, controllers.CreateTag)
	router.GET("/", controllers.GetAllTags)

	router.GET("/:id", controllers.GetTagById)
	router.PATCH("/:id", middlewares.CheckAuth, controllers.UpdateTag)
	router.DELETE("/:id", middlewares.CheckAuth, controllers.DeleteTagById)

}
