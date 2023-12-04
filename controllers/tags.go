package controllers

import (
	"net/http"

	"github.com/Kunniii/gocms/internal"
	"github.com/Kunniii/gocms/models"
	"github.com/gin-gonic/gin"
)

func CreateTag(context *gin.Context) {
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

	tag := models.Tag{Name: reqBody.Name}

	result := internal.DB.Create(&tag)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Could not create Tag!",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"OK":   true,
		"data": tag,
	})

}
func GetAllTags(context *gin.Context) {
	var tags []models.Tag
	internal.DB.Find(&tags)

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
