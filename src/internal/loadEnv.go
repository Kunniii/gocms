package internal

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
		log.Println("Using current environment variables!")
	} else {
		log.Println("Loaded environment variables!")
	}
}
