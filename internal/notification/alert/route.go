package alert

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	alerts := r.Group("/alerts", middleware.AuthMiddleware())
	{
		alerts.GET("", h.FindAll)
		alerts.GET("/:id", h.FindByID)
		alerts.PUT("/:id/read", h.MarkRead)
	}
}
