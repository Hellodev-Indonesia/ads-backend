package ads

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.GET("/meta/ads", middleware.AuthMiddleware(), middleware.RequirePermission("meta.campaign.view"), h.GetAds)
	r.GET("/meta/creatives/:id", middleware.AuthMiddleware(), middleware.RequirePermission("meta.campaign.view"), h.GetCreative)
}
