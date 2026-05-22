package dashboard

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/meta/dashboard", middleware.AuthMiddleware())
	{
		g.GET("/campaigns", middleware.RequirePermission("meta.campaign.view"), h.GetCampaignDashboard)
	}
}
