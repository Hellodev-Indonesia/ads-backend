package auth

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", h.Login)
		authGroup.POST("/logout", middleware.AuthMiddleware(), h.Logout)
	}
}
