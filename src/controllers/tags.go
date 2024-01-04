package controllers

import (
	"net/http"

	"github.com/Kunniii/gocms/apiModels"
	"github.com/Kunniii/gocms/internal"
	"github.com/Kunniii/gocms/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CreateTag(context *gin.Context) {

	authToken := context.GetString("auth-token")
	token, _, _ := internal.VerifyToken(authToken)
	userClaims := token.Claims.(jwt.MapClaims)
	if roleId := uint(userClaims["RoleID"].(float64)); roleId < 2 {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"OK":      false,
			"message": "Forbidden",
		})
		return
	}

	var reqBody struct {
		Name string
	}

	if err := context.Bind(&reqBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Make sure to put the JSON key as String and no trailing commas",
		})
		return
	}

	var tag models.Tag
	result := internal.DB.First(&tag, "name = ?", reqBody.Name)

	if result.Error == nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Tag already exists!",
		})
		return
	}

	tag = models.Tag{Name: reqBody.Name}

	internal.DB.Create(&tag)

	context.JSON(http.StatusOK, gin.H{
		"OK":   true,
		"data": tag,
	})

}

func GetAllTags(context *gin.Context) {
	var tags []apiModels.Tag
	internal.DB.Model(&models.Tag{}).Find(&tags)

	context.JSON(http.StatusOK, gin.H{
		"OK":   true,
		"data": tags,
	})

}

func GetTagById(context *gin.Context) {
	id := context.Param("id")

	var tag models.Tag

	if result := internal.DB.First(&tag, id); result.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"OK":      false,
			"message": "Not found!",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"OK":   true,
			"data": tag,
		})
	}
}

func UpdateTag(context *gin.Context) {

	authToken := context.GetString("auth-token")
	token, _, _ := internal.VerifyToken(authToken)
	userClaims := token.Claims.(jwt.MapClaims)
	if roleId := uint(userClaims["RoleID"].(float64)); roleId < 2 {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"OK":      false,
			"message": "Forbidden",
		})
		return
	}

	id := context.Param("id")

	var reqBody struct {
		Name string
	}

	if err := context.Bind(&reqBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Make sure to put the JSON key as String and no trailing commas",
		})
		return
	}

	var tag models.Tag
	if result := internal.DB.First(&tag, id); result.Error != nil {

		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "ID not found",
		})
		return
	}

	internal.DB.Model(&tag).Updates(models.Tag{
		Name: reqBody.Name,
	})

	context.JSON(http.StatusOK, gin.H{
		"OK":   true,
		"data": tag,
	})
}

func DeleteTagById(context *gin.Context) {
	id := context.Param("id")

	if result := internal.DB.Delete(&models.Tag{}, id); result.Error != nil {
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
