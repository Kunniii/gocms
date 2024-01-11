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
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

func CreateAdmin(user *models.User) {
	hashByte, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Cannot generate Admin password!")
	}

	user.Password = string(hashByte)

	if result := internal.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(user); result.Error != nil {
		err := result.Error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			log.Println("Admin email already exists")
			return
		} else {
			log.Fatal("Cannot create admin user!")
			return
		}
	}

}

func Register(context *gin.Context) {
	var reqBody struct {
		UserName string
		Email    string
		Password string
	}

	if context.Bind(&reqBody) != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Make sure to put the JSON key as String and no trailing commas",
		})
	}

	if reqBody.UserName == "" || reqBody.Email == "" || reqBody.Password == "" {
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
		RoleID:   itypes.Roles[0].ID,
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

func Login(context *gin.Context) {
	var reqBody struct {
		Email    string
		Password string
	}

	if err := context.Bind(&reqBody); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Make sure to put the JSON key as String and no trailing commas",
		})
	}
	var user models.User
	result := internal.DB.First(&user, "email = ?", reqBody.Email)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Invalid credential!",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "Invalid credential!",
		})
		return
	}

	tokenString, err := internal.CreateToken(itypes.UserClaims{
		UserID:   user.ID,
		RoleID:   user.RoleID,
		Email:    user.Email,
		UserName: user.UserName,
	})

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"OK":      false,
			"message": "Cannot generate token!",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"OK":    true,
			"Token": "Bearer " + tokenString,
		})
	}
}

func Verify(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"OK": true,
	})
}

func Me(context *gin.Context) {

	authToken := context.GetString("auth-token")
	token, _, _ := internal.VerifyToken(authToken)
	userClaims := token.Claims.(jwt.MapClaims)
	userID := uint(userClaims["UserID"].(float64))

	var user apiModels.User
	internal.DB.Model(&models.User{}).First(&user, userID)

	context.JSON(http.StatusOK, gin.H{
		"OK":   true,
		"data": user,
	})

}
