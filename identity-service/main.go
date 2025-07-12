package main

import (
	"identity-service/config"
	"identity-service/models"
	"identity-service/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	config.DB.AutoMigrate(&models.User{}, &models.RefreshToken{})
	PORT := os.Getenv("PORT")
	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run(":" + PORT)
}
