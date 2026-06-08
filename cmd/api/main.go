package main

import (
	"log"

	"github.com/alex/ads_backend/config"
	"github.com/alex/ads_backend/routes"
	"github.com/gin-gonic/gin"
)

// @title           Ads Backend API
// @version         1.0
// @description     ERP Ready Backend for Meta Ads & Business Management
// @host            localhost:8888
// @BasePath        /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Tuliskan 'Bearer ' diikuti dengan token JWT Anda.
func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize configurations
	config.InitDB()
	config.InitRedis()
	config.InitMeta()
	config.InitCentrifugo()

	// Setup Gin
	router := gin.Default()

	// Register routes (Pure Gin)
	routes.RegisterApiRoutes(router)

	// Start server
	port := config.GetEnv("APP_PORT", "8080")
	log.Printf("Server starting on port %s", port)
	log.Printf("Documentation available at: http://localhost:%s/docs", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
