package reporting

import (
	"github.com/gin-gonic/gin"
	
	"github.com/alex/ads_backend/internal/reporting/dashboard"
	"gorm.io/gorm"
)

func RegisterReportingRoutes(router *gin.RouterGroup, db *gorm.DB) {
	reportingRoutes := router.Group("/reporting")
	{
		dashboard.RegisterRoutes(reportingRoutes, db)
	}
}
