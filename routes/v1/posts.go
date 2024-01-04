package v1

import (
	"net/http"

	"github.com/Kunniii/gocms/controllers"
	"github.com/Kunniii/gocms/middlewares"
	"github.com/gin-gonic/gin"
)

func UsePostsRouter(router *gin.RouterGroup) {
	router.GET("/ready", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"OK": true,
		})
	})

	router.POST("/", middlewares.CheckAuth, controllers.CreatePost)
	router.GET("/", controllers.GetAllPosts)

	router.GET("/:id", controllers.GetPostById)
	router.PATCH("/:id", middlewares.CheckAuth, controllers.UpdatePost)
	router.DELETE("/:id", middlewares.CheckAuth, controllers.DeletePostById)

	router.POST("/:id/comments", middlewares.CheckAuth, controllers.AddComment)

	router.GET("/:id/comments", controllers.GetComments)
	router.GET("/:id/comments/:offset", controllers.GetComments)

	router.GET("/:id/tags", controllers.GetTags)

	router.GET("/:id/likes", controllers.GetLikes)
	router.GET("/:id/likes/:offset", controllers.GetLikes)
	router.POST("/:id/likes", middlewares.CheckAuth, controllers.LikePost)

}
