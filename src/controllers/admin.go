package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Kunniii/gocms/apiModels"
	"github.com/Kunniii/gocms/internal"
	itypes "github.com/Kunniii/gocms/internal/types"
	"github.com/Kunniii/gocms/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
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

	if reqBody.Email == "" ||
		reqBody.RoleID > 4 ||
		reqBody.UserName == "" ||
		reqBody.Password == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Missing data!",
		})
		return
	}

	// hash user's password
	hashByte, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"OK":      false,
			"message": "Password hashing error!",
		})
		return
	}

	// by default, user is registered as normal user
	// with id = 0
	user := models.User{
		UserName: reqBody.UserName,
		Email:    reqBody.Email,
		Password: string(hashByte),
		RoleID:   itypes.Roles[reqBody.RoleID].ID,
	}

	if result := internal.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&user); result.Error != nil {
		err := result.Error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			context.JSON(http.StatusBadRequest, gin.H{
				"OK":      false,
				"message": "Email already exists!",
			})
			return
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{
				"OK":      false,
				"message": "Cannot create user!",
			})
			return
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"OK": true,
	})
}

func GetUsers(context *gin.Context) {
	var users []apiModels.User

	internal.DB.Model(&models.User{}).Find(&users)

	context.JSON(http.StatusOK, gin.H{
		"OK":   true,
		"data": users,
	})
}

func DeleteUserByID(context *gin.Context) {
	var id = context.Param("id")

	if err := internal.DB.Delete(&models.User{}, id).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": err,
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"OK": true,
		})
	}
}
