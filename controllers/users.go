package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Kunniii/gocms/internal"
	"github.com/Kunniii/gocms/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func Register(context *gin.Context) {
	var reqBody struct {
		UserName string
		Email    string
		Password string
	}

	if context.Bind(&reqBody) != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":  false,
			"msg": "Make sure to put the JSON key as String",
		})
		return
	}

	// hash user's password
	hashByte, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"OK":  false,
			"msg": "Password hashing error!",
		})
		return
	}

	user := models.User{UserName: reqBody.UserName, Email: reqBody.Email, Password: string(hashByte)}

	if result := internal.DB.Create(&user); result.Error != nil {
		err := result.Error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			context.JSON(http.StatusBadRequest, gin.H{
				"OK":  false,
				"msg": "Email already exists!",
			})
			return
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{
				"OK":  false,
				"msg": "Cannot create user!",
			})
			return
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"OK": true,
	})
}

func Login(context *gin.Context) {
	var reqBody struct {
		Email    string
		Password string
	}

	if err := context.Bind(&reqBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":  false,
			"msg": "Make sure to put JSON key as String!",
		})
		return
	}

	var user models.User
	result := internal.DB.First(&user, "email = ?", reqBody.Email)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":  false,
			"msg": "Invalid credential!",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":  false,
			"msg": "Invalid credential!",
		})
		return
	}

	tokenString, err := internal.CreateToken()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"OK":  false,
			"msg": "Cannot generate token!",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"OK":    true,
			"token": "Bearer " + tokenString,
		})
	}
}

func Validate(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"OK": true,
	})
}
