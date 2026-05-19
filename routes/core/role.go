package core

import (
	"github.com/alex/ads_backend/internal/core/role"
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoleRoutes(r *gin.RouterGroup, h *role.Handler) {
	roles := r.Group("/roles")
	roles.Use(middleware.AuthMiddleware())
	{
		roles.GET("", h.FindAll)
		roles.GET("/:id", h.FindByID)
		roles.POST("", h.Create)
		roles.PUT("/:id", h.Update)
		roles.DELETE("/:id", h.Delete)
		roles.POST("/:id/permissions", h.AssignPermissions)
	}
}
