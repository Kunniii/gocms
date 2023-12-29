package controllers

import (
	"log"
	"net/http"

	"github.com/Kunniii/gocms/internal"
	"github.com/Kunniii/gocms/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CreatePost(context *gin.Context) {

	var reqBody struct {
		Title string
		Body  string
		Tags  []uint
	}

	if err := context.Bind(&reqBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Make sure to put the JSON key as String and no trailing commas",
		})
		return
	}

	authToken := context.GetString("auth-token")
	userClaims := internal.GetClaims(authToken)
	userId := uint(userClaims["UserID"].(float64))

	if roleId := uint(userClaims["RoleID"].(float64)); roleId < 1 {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"OK":      false,
			"message": "Forbidden",
		})
		return
	}

	var tags []models.Tag
	internal.DB.Find(&tags, reqBody.Tags)

	var user models.User
	internal.DB.First(&user, userId)

	post := models.Post{
		Title:  reqBody.Title,
		Body:   reqBody.Body,
		UserID: userId,
		Tags:   tags,
	}

	result := internal.DB.Create(&post)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Could not create Post!",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"OK":   true,
		"data": post,
	})
}

func GetAllPosts(context *gin.Context) {
	var posts []models.Post

	internal.DB.Select([]string{"id", "updated_at", "user_id", "title"}).Find(&posts)

	context.JSON(http.StatusOK, gin.H{
		"OK":   true,
		"data": posts,
	})
}

func GetPostById(context *gin.Context) {
	id := context.Param("id")

	var post models.Post
	if err := internal.DB.Model(&models.Post{}).Preload("Tags").Preload("Comments").First(&post, id).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"OK":      false,
			"message": "Not found!",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"OK":   true,
			"data": post,
		})
	}
}

func UpdatePost(context *gin.Context) {
	id := context.Param("id")

	var reqBody struct {
		Title string
		Body  string
		Tags  []uint
	}

	if err := context.Bind(&reqBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Make sure to put the JSON key as String and no trailing commas",
		})
		return
	}

	authToken := context.GetString("auth-token")
	token, _, _ := internal.VerifyToken(authToken)
	userClaims := token.Claims.(jwt.MapClaims)
	if roleId := uint(userClaims["RoleID"].(float64)); roleId < 1 {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"OK":      false,
			"message": "Forbidden",
		})
		return
	}

	var tags []*models.Tag
	if len(reqBody.Tags) > 0 {
		internal.DB.Find(&tags, reqBody.Tags)
	}

	var post models.Post
	if result := internal.DB.First(&post, id); result.Error != nil {

		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "ID not found",
		})
		return
	}

	internal.DB.Model(&post).Updates(models.Post{
		Title: reqBody.Title,
		Body:  reqBody.Body,
	})

	if err := internal.DB.Model(&post).Association("Tags").Replace(tags); err != nil {
		log.Println(err)
	}

	context.JSON(http.StatusOK, gin.H{
		"OK":   true,
		"data": post,
	})
}

func DeletePostById(context *gin.Context) {
	id := context.Param("id")

	if result := internal.DB.Delete(&models.Post{}, id); result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": result.Error,
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"OK": true,
		})
	}
}

func AddComment(context *gin.Context) {
	postID := context.Param("id")

	var reqBody struct {
		Body string
	}

	if err := context.Bind(&reqBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Make sure to put the JSON key as String and no trailing commas",
		})
		return
	}

	var post models.Post
	if err := internal.DB.Model(&models.Post{}).Preload("Comments").First(&post, postID).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"OK":      false,
			"message": "Post not found!",
		})
		return
	}

	// get userID from jwt
	authToken := context.GetString("auth-token")
	userClaims := internal.GetClaims(authToken)
	userId := uint(userClaims["UserID"].(float64))

	// create comment
	var comment = models.Comment{
		UserID: userId,
		PostID: post.ID,
		Body:   reqBody.Body,
	}

	post.Comments = append(post.Comments, comment)

	if err := internal.DB.Save(&post).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"OK":      false,
			"message": err,
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"OK":   true,
			"Post": post,
		})
	}

}

func GetComment(context *gin.Context) {}
