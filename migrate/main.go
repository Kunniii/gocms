package main

import (
	"log"
	"os"

	"github.com/Kunniii/gocms/controllers"
	"github.com/Kunniii/gocms/internal"
	itypes "github.com/Kunniii/gocms/internal/types"
	"github.com/Kunniii/gocms/models"
)

func init() {
	internal.LoadEnv()
	internal.ConnectDB()
}

func main() {
	var dbModels = []interface{}{
		&models.Comment{},
		&models.Post{},
		&models.Tag{},
		&models.User{},
	}

	if err := internal.DB.AutoMigrate(dbModels...); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database migration successfully!")
	}

	adminEmail := os.Getenv("DEFAULT_ADMIN_EMAIL")
	adminUsername := os.Getenv("DEFAULT_ADMIN_USERNAME")
	adminPassword := os.Getenv("DEFAULT_ADMIN_PASSWORD")

	adminUser := models.User{
		UserName: adminUsername,
		Email:    adminEmail,
		Password: adminPassword,
		RoleID:   itypes.Roles[4].ID,
	}
	controllers.CreateAdmin(&adminUser)
}
