package adset

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.GET("/meta/adsets", middleware.AuthMiddleware(), middleware.RequirePermission("meta.adset.view"), h.GetAdSets)
	r.GET("/meta/brands/:brand_id/adsets", middleware.AuthMiddleware(), middleware.RequirePermission("meta.adset.view"), h.GetAdSetsByBrand)
}
