package meta

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	metaGroup := r.Group("/meta")
	metaGroup.Use(middleware.AuthMiddleware())
	{
		metaGroup.GET("/ad-accounts", middleware.RequirePermission("meta.campaign.view"), h.GetAdAccounts)
		metaGroup.GET("/campaigns", middleware.RequirePermission("meta.campaign.view"), h.GetCampaigns)
		metaGroup.GET("/adsets", middleware.RequirePermission("meta.campaign.view"), h.GetAdSets)
		metaGroup.GET("/ads", middleware.RequirePermission("meta.campaign.view"), h.GetAds)
		metaGroup.GET("/insights", middleware.RequirePermission("meta.campaign.view"), h.GetInsights)
	}
}
