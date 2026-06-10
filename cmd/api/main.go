package main

import (
	"context"
	"log"

	"github.com/alex/ads_backend/config"
	metasync "github.com/alex/ads_backend/internal/meta/sync"
	"github.com/alex/ads_backend/routes"
	"github.com/gin-contrib/cors"
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

	// Clean up orphaned sync batches
	cleanupOrphanedSyncBatches()

	// Setup Gin
	router := gin.Default()

	// CORS Middleware
	configCORS := cors.DefaultConfig()
	configCORS.AllowAllOrigins = true
	configCORS.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(configCORS))

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

func cleanupOrphanedSyncBatches() {
	importContext := context.Background()
	repo := metasync.NewRepository(config.DB)
	service := metasync.NewService(repo)
	if err := service.CleanupOrphanedBatches(importContext); err != nil {
		log.Printf("Failed to clean up orphaned sync batches: %v", err)
	} else {
		log.Println("Successfully cleaned up orphaned sync batches")
	}
}
