package role

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
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
