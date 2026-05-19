package core

import (
	"github.com/alex/ads_backend/internal/core/user"
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup, h *user.Handler) {
	users := r.Group("/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("", h.FindAll)
		users.GET("/:id", h.FindByID)
		users.POST("", h.Create)
		users.PUT("/:id", h.Update)
		users.DELETE("/:id", h.Delete)
	}
}
