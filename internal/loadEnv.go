package internal

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Loaded environment variables!")
	}
}
