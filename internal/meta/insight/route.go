package insight

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.GET("/meta/insights", middleware.AuthMiddleware(), middleware.RequirePermission("meta.campaign.view"), h.GetInsights)
}
