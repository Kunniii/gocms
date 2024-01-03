package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Kunniii/gocms/apiModels"
	"github.com/Kunniii/gocms/internal"
	"github.com/Kunniii/gocms/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// #region Post

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
	if err := internal.DB.Model(&models.Post{}).First(&post, id).Error; err != nil {
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

// #endregion

// #region Comments

func AddComment(context *gin.Context) {
	postID := context.Param("id")

	var reqBody struct {
		Body string
	}

	if err := context.Bind(&reqBody); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Make sure to put the JSON key as String and no trailing commas",
		})
		return
	}

	var post models.Post
	if err := internal.DB.Model(&models.Post{}).Preload("Comments").First(&post, postID).Error; err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
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
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"OK":      false,
			"message": err,
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"OK":   true,
			"data": comment,
		})
	}
}

func GetComment(context *gin.Context) {
	postID := context.Param("id")
	offset := context.Param("offset")
	var pageNumber int = 0
	var err error

	if offset != "" {
		pageNumber, err = strconv.Atoi(offset)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"OK":      false,
				"message": "Cannot get that offset!",
			})
			return
		} else {
			// starts at 1
			pageNumber -= 1
		}
	}

	var comments []apiModels.Comment
	internal.DB.Model(&models.Comment{}).Where("post_id = ?", postID).Order("id desc").Limit(5).Offset(5 * pageNumber).Find(&comments)

	context.JSON(http.StatusOK, gin.H{
		"OK":   true,
		"data": comments,
	})
}

// #endregion
