package ads

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.GET("/meta/ads", middleware.AuthMiddleware(), middleware.RequirePermission("meta.ads.view"), h.GetAds)
	r.GET("/meta/brands/:brand_id/ads", middleware.AuthMiddleware(), middleware.RequirePermission("meta.ads.view"), h.GetAdsByBrand)
	r.GET("/meta/creatives/:id", middleware.AuthMiddleware(), middleware.RequirePermission("meta.ads.view"), h.GetCreative)
}
