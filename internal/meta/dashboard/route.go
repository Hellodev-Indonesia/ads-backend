package dashboard

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/meta/dashboard", middleware.AuthMiddleware())
	{
		g.GET("/campaigns", middleware.RequirePermission("meta.campaign.view"), h.GetCampaignDashboard)
		g.GET("/adsets", middleware.RequirePermission("meta.campaign.view"), h.GetAdSetDashboard)
		g.GET("/ads", middleware.RequirePermission("meta.campaign.view"), h.GetAdDashboard)
	}
}
