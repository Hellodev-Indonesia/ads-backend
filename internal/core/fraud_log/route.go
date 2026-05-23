package fraud_log

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	logs := r.Group("/fraud-logs", middleware.AuthMiddleware())
	{
		logs.GET("", h.FindAll)
		logs.GET("/:id", h.FindByID)
		logs.PUT("/:id/resolve", h.Resolve)
	}
}
