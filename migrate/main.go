package main

import (
	"log"

	"github.com/Kunniii/gocms/internal"
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
}
