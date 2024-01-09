package controllers

import (
	"log"
	"net/http"

	"github.com/Kunniii/gocms/apiModels"
	"github.com/Kunniii/gocms/internal"
	"github.com/Kunniii/gocms/models"
	"github.com/gin-gonic/gin"
)

// assume the user is authorized!

func AddUser(context *gin.Context) {
	var reqBody struct {
		RoleID   uint
		Email    string
		UserName string
		Password string
	}

	if err := context.Bind(&reqBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Make sure to put the JSON key as String and no trailing commas",
		})
		return
	}

	var user = models.User{
		UserName: reqBody.UserName,
		Email:    reqBody.Email,
		RoleID:   reqBody.RoleID,
		Password: reqBody.Password,
	}

	if err := internal.DB.Model(&models.User{}).Create(&user).Error; err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": err,
		})
		return
	}

	var apiUser apiModels.User

	if err := internal.DB.Model(&models.User{}).First(&apiUser, user.ID); err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": err,
		})
		return
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":   true,
			"data": apiUser,
		})
	}
}

func UpdateUser(context *gin.Context) {

}

func DeleteUser(context *gin.Context) {

}
