package dashboard

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/alex/ads_backend/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {
	repo := NewDashboardRepository(db)
	service := NewDashboardService(repo)
	handler := NewDashboardHandler(service)

	dashboardRoutes := router.Group("/dashboard")
	// Require permission reporting.dashboard.view to access these endpoints
	dashboardRoutes.Use(middleware.AuthMiddleware())
	dashboardRoutes.Use(middleware.RequirePermission("reporting.dashboard.view"))
	{
		dashboardRoutes.GET("/summary", handler.GetSummary)
		dashboardRoutes.GET("/brands", handler.GetBrandsMonitoring)
		dashboardRoutes.GET("/activities", handler.GetRecentActivities)
	}
}
