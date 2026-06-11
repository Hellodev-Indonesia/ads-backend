package user

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	users := r.Group("/core/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("", h.FindAll)
		users.GET("/:id", h.FindByID)
		users.POST("", h.Create)
		users.PUT("/:id", h.Update)
		users.DELETE("/:id", h.Delete)
	}
}
