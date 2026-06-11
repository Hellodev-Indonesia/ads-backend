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
		brands.GET("/dashboard", middleware.RequirePermission("core.brand.view"), h.GetBrandDashboard)
		brands.GET("/:slug", middleware.RequirePermission("core.brand.view"), h.FindBySlug)
		brands.POST("", middleware.RequirePermission("core.brand.create"), h.Create)
		brands.PUT("/:slug", middleware.RequirePermission("core.brand.update"), h.Update)
		brands.DELETE("/:slug", middleware.RequirePermission("core.brand.delete"), h.Delete)
	}
}
