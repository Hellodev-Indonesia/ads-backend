package core

import (
	"github.com/alex/ads_backend/internal/core/permission"
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterPermissionRoutes(r *gin.RouterGroup, h *permission.Handler) {
	perms := r.Group("/permissions")
	perms.Use(middleware.AuthMiddleware())
	{
		perms.GET("", h.FindAll)
		perms.GET("/:id", h.FindByID)
		perms.POST("", h.Create)
		perms.PUT("/:id", h.Update)
		perms.DELETE("/:id", h.Delete)
	}
}
