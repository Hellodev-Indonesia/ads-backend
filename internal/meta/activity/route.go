package activity

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.GET("/meta/brands/:brand_id/activities", middleware.AuthMiddleware(), middleware.RequirePermission("meta.activity.view"), h.GetActivitiesByBrand)
}
