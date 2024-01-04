package controllers

import (
	"net/http"

	"github.com/Kunniii/gocms/internal"
	"github.com/Kunniii/gocms/models"
	"github.com/gin-gonic/gin"
)

func DeleteComment(context *gin.Context) {
	commentID := context.Param("id")

	if commentID == "" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Need comment id",
		})
		return
	}

	// find the comment with id
	var comment models.Comment
	internal.DB.First(&comment, commentID)

	// get userID from jwt
	authToken := context.GetString("auth-token")
	userClaims := internal.GetClaims(authToken)
	userID := uint(userClaims["UserID"].(float64))
	// compare with userID

	if comment.UserID != userID {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"OK":      false,
			"message": "Not your Comment!",
		})
	} else {
		internal.DB.Delete(&comment)
		context.JSON(http.StatusOK, gin.H{
			"OK":   true,
			"data": comment,
		})
	}
}

func UpdateComment(context *gin.Context) {
	commentID := context.Param("id")

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

	if commentID == "" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Need comment id",
		})
		return
	}

	// find the comment with id
	var comment models.Comment
	internal.DB.First(&comment, commentID)

	// get userID from jwt
	authToken := context.GetString("auth-token")
	userClaims := internal.GetClaims(authToken)
	userID := uint(userClaims["UserID"].(float64))
	// compare with userID

	if comment.UserID != userID {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"OK":      false,
			"message": "Not your Comment!",
		})
	} else {
		internal.DB.Model(&comment).Updates(models.Comment{
			Body: reqBody.Body,
		})
		context.JSON(http.StatusOK, gin.H{
			"OK":   true,
			"data": comment,
		})
	}
}
