package brand_whitelist_rule

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	rules := r.Group("/brands/:brand_id/whitelist-rules", middleware.AuthMiddleware())
	{
		rules.POST("", h.Create)
		rules.GET("", h.FindAll)
		rules.GET("/:id", h.FindByID)
		rules.PUT("/:id", h.Update)
		rules.DELETE("/:id", h.Delete)
		rules.POST("/check-url", h.CheckURL)
	}
}
