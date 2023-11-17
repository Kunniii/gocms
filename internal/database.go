package internal

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseInfo struct {
	host     string
	user     string
	password string
	dbname   string
	port     string
	sslmode  string
	TimeZone string
}

func (dbInfo *DatabaseInfo) toString() string {
	return fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v", dbInfo.host, dbInfo.user, dbInfo.password, dbInfo.dbname, dbInfo.port, dbInfo.sslmode, dbInfo.TimeZone)
}

var DB *gorm.DB

func ConnectDB() {
	var dbInfo DatabaseInfo = DatabaseInfo{
		host:     os.Getenv("DB_HOST"),
		user:     os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASS"),
		dbname:   os.Getenv("DB_NAME"),
		port:     os.Getenv("DB_PORT"),
		sslmode:  os.Getenv("DB_SSLMODE"),
		TimeZone: os.Getenv("DB_TIMEZONE"),
	}
	var err error
	dsn := dbInfo.toString()
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Cannot connect to database!")
	} else {
		log.Println("Database connection is ready!")
	}
}
