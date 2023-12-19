package controllers

import (
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
		Tags  []int
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
	userId := uint(userClaims["UserID"].(float64))

	// if roleId := uint(userClaims["RoleID"].(float64)); roleId < 1 {
	// 	context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
	// 		"OK":      false,
	// 		"message": "Forbidden",
	// 	})
	// 	return
	// }

	var tags []*models.Tag
	internal.DB.Find(&tags, reqBody.Tags)

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

	internal.DB.Model(&models.Post{}).Preload("Tags").Preload("User").Find(&posts)

	// internal.DB.Find(&posts)

	context.JSON(http.StatusOK, gin.H{
		"OK":   true,
		"data": posts,
	})

}

func GetPostById(context *gin.Context) {
	id := context.Param("id")

	var post models.Post
	if err := internal.DB.Model(&models.Post{}).Preload("Tags").First(&post, id).Error; err != nil {
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
	}

	if err := context.Bind(&reqBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Make sure to put the JSON key as String and no trailing commas",
		})
		return
	}

	// authToken := context.GetString("auth-token")
	// token, _, _ := internal.VerifyToken(authToken)
	// userClaims := token.Claims.(jwt.MapClaims)
	// userId := uint(userClaims["UserID"].(float64))

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
