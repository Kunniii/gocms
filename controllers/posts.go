package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Kunniii/gocms/apiModels"
	"github.com/Kunniii/gocms/internal"
	"github.com/Kunniii/gocms/models"
	"github.com/gin-gonic/gin"
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

	if reqBody.Title == "" && reqBody.Body == "" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Empty Post",
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
	if len(reqBody.Tags) > 0 {
		internal.DB.Find(&tags, reqBody.Tags)
	}

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

	var apiPost apiModels.Post
	internal.DB.Model(&models.Post{}).Where("id = ?", post.ID).Find(&apiPost)

	context.JSON(http.StatusOK, gin.H{
		"OK":   true,
		"data": apiPost,
	})
}

func GetAllPosts(context *gin.Context) {
	var posts []apiModels.Post

	internal.DB.Model(&models.Post{}).Find(&posts)

	context.JSON(http.StatusOK, gin.H{
		"OK":   true,
		"data": posts,
	})
}

func GetPostById(context *gin.Context) {
	id := context.Param("id")

	var post apiModels.Post
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
	userClaims := internal.GetClaims(authToken)
	if roleId := uint(userClaims["RoleID"].(float64)); roleId < 1 {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"OK":      false,
			"message": "Forbidden",
		})
		return
	}

	var post models.Post
	if result := internal.DB.First(&post, id); result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Not found",
		})
		return
	}

	if userID := uint(userClaims["UserID"].(float64)); userID != post.UserID {
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

	internal.DB.Model(&post).Updates(models.Post{
		Title: reqBody.Title,
		Body:  reqBody.Body,
	})

	if err := internal.DB.Model(&post).Association("Tags").Replace(tags); err != nil {
		log.Println(err)
	}

	var apiPost apiModels.Post
	internal.DB.Model(&models.Post{}).Where("id = ?", post.ID).Find(&apiPost)

	context.JSON(http.StatusOK, gin.H{
		"OK":   true,
		"data": apiPost,
	})
}

func DeletePostById(context *gin.Context) {
	id := context.Param("id")

	// TODO: check if post belongs to user

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

func GetComments(context *gin.Context) {
	postID := context.Param("id")
	offset := context.Param("offset")
	var pageNumber int = 0
	var err error

	if offset != "" {
		pageNumber, err = strconv.Atoi(offset)
		if err != nil || pageNumber < 0 {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"OK":      false,
				"message": "Cannot get that offset!",
			})
			return
		} else {
			// this is for when user input = 1
			// which means instead of start at 0
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

// #region Like

func LikePost(context *gin.Context) {
	postID := context.Param("id")

	authToken := context.GetString("auth-token")
	userClaims := internal.GetClaims(authToken)
	userId := uint(userClaims["UserID"].(float64))

	var like models.Like

	internal.DB.Where("user_id = ? AND post_id = ?", userId, postID).Find(&like)

	if like.ID == 0 {
		var post models.Post
		internal.DB.Find(&post, postID)
		like.UserID = userId
		like.PostID = post.ID
		post.Likes = append(post.Likes, like)
		internal.DB.Save(&post)

		context.JSON(http.StatusOK, gin.H{
			"OK":      true,
			"message": "Add like",
		})
	} else {
		internal.DB.Where("id = ?", like.ID).Delete(&like)
		context.JSON(http.StatusOK, gin.H{
			"OK":      true,
			"message": "Remove like",
		})
	}
}

func GetLikes(context *gin.Context) {
	postID := context.Param("id")
	offset := context.Param("offset")
	var pageNumber int = 0
	var err error

	if offset != "" {
		pageNumber, err = strconv.Atoi(offset)
		if err != nil || pageNumber < 0 {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"OK":      false,
				"message": "Cannot get that offset!",
			})
			return
		} else {
			// this is for when user input = 1
			// which means instead of start at 0
			// starts at 1
			pageNumber -= 1
		}
	}

	var likes []apiModels.Like
	internal.DB.Model(&models.Like{}).Where("post_id = ?", postID).Order("id desc").Limit(5).Offset(5 * pageNumber).Find(&likes)

	var total int64

	internal.DB.Model(&models.Like{}).Where("post_id = ?", postID).Count(&total)

	context.JSON(http.StatusOK, gin.H{
		"OK": true,
		"data": gin.H{
			"total": total,
			"likes": likes,
		},
	})
}

// #endregion

func GetTags(context *gin.Context) {
	postID := context.Param("id")

	var post models.Post

	if err := internal.DB.First(&post, postID).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Not found",
		})
	} else {
		// find all id that has post_id = postID, and select tag_id
		var tagIDs []uint
		internal.DB.Table("posts_tags").Select("tag_id").Where("post_id = ?", postID).Find(&tagIDs)

		if len(tagIDs) == 0 {
			context.JSON(http.StatusOK, gin.H{
				"OK":   true,
				"data": tagIDs, //because it is 0, return will be an [] array
			})
		} else {
			var tags []apiModels.Tag
			internal.DB.Model(&models.Tag{}).Order("id desc").Find(&tags, tagIDs)

			context.JSON(http.StatusOK, gin.H{
				"OK":   true,
				"data": tags,
			})
		}
	}
}
