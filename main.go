package main

import (
	"log"

	"github.com/Kunniii/gocms/internal"
	v1 "github.com/Kunniii/gocms/routes/v1"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	internal.LoadEnv()
	internal.ConnectDB()
}

func main() {
	api := v1.NewRouter()

	log.Fatal(api.Run())
}
