package campaign

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.GET("/meta/campaigns", middleware.AuthMiddleware(), middleware.RequirePermission("meta.campaign.view"), h.GetCampaigns)
	r.GET("/meta/campaigns/dashboard", middleware.AuthMiddleware(), middleware.RequirePermission("meta.campaign.view"), h.GetCampaignDashboard)
	r.GET("/meta/brands/:brand_id/campaigns", middleware.AuthMiddleware(), middleware.RequirePermission("meta.campaign.view"), h.GetCampaignsByBrand)
	r.GET("/meta/brands/:brand_id/campaigns/summary", middleware.AuthMiddleware(), middleware.RequirePermission("meta.campaign.view"), h.GetCampaignSummaryByBrand)
	r.GET("/meta/brands/:brand_id/campaigns/list", middleware.AuthMiddleware(), middleware.RequirePermission("meta.campaign.view"), h.GetCampaignListByBrand)
}
