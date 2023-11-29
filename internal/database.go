package internal

import (
	"log"
	"os"

	itypes "github.com/Kunniii/gocms/internal/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var dbInfo itypes.DatabaseInfo = itypes.DatabaseInfo{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		TimeZone: os.Getenv("DB_TIMEZONE"),
	}
	var err error
	dsn := dbInfo.ToString()
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Cannot connect to database!")
	} else {
		log.Println("Database connection is ready!")
	}
}
