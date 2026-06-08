package brand

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	brands := r.Group("/core/brands")
	brands.Use(middleware.AuthMiddleware())
	{
		brands.GET("", middleware.RequirePermission("core.brand.view"), h.FindAll)
		brands.GET("/:id", middleware.RequirePermission("core.brand.view"), h.FindByID)
		brands.POST("", middleware.RequirePermission("core.brand.create"), h.Create)
		brands.PUT("/:id", middleware.RequirePermission("core.brand.update"), h.Update)
		brands.DELETE("/:id", middleware.RequirePermission("core.brand.delete"), h.Delete)
	}
}
