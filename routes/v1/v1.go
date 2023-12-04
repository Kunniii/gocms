package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")
	v1.GET("/ready", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"OK": true,
		})
	})
	UsePostsRouter(v1.Group("/posts"))
	UseUsersRouter(v1.Group("/users"))
	UseTagRouter(v1.Group("/tags"))
	return router
}
