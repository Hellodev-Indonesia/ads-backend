package insight

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	insights := r.Group("/meta/insights")
	{
		insights.GET("/campaign", middleware.AuthMiddleware(), middleware.RequirePermission("meta.insight.view"), h.GetCampaignInsights)
		insights.GET("/ad", middleware.AuthMiddleware(), middleware.RequirePermission("meta.insight.view"), h.GetAdInsights)
	}
}
