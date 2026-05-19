package core

import (
	"github.com/alex/ads_backend/internal/core/auth"
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.RouterGroup, h *auth.Handler) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", h.Login)
		authGroup.POST("/logout", middleware.AuthMiddleware(), h.Logout)
	}
}
