package controllers

import (
	"log"

	"github.com/Kunniii/gocms/internal"
	"github.com/Kunniii/gocms/models"
)

func CreateComment(userID uint, postID uint, body string) (*models.Comment, bool) {
	var comment = models.Comment{
		UserID: userID,
		PostID: postID,
		Body:   body,
	}

	if err := internal.DB.Create(&comment).Error; err != nil {
		log.Println(err)
		return nil, false
	}

	return &comment, true
}

func DeleteComment(commentID uint, userID uint) bool {
	// find the comment with id
	// compare with userID
	return true
}

func UpdateComment(commentID uint, body string) (*models.Comment, bool) {
	// find the comment with id
	// compare with userID
	return nil, false
}

func GetCommentWithID(commentID uint) (*models.Comment, bool) {
	return nil, false
}
