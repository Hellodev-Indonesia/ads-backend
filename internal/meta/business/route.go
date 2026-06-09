package business

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	businesses := r.Group("/businesses")
	{
		businesses.GET("", middleware.RequirePermission("meta.business.view"), h.GetBusinesses)
	}
}
