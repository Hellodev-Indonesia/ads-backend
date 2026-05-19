package middleware

import (
	"strings"

	"github.com/alex/ads_backend/pkg/response"
	"github.com/alex/ads_backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware memvalidasi token JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "Authorization header must be Bearer token")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set data user ke context Gin
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("roles", claims.Roles)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}

// RequirePermission mengecek apakah user punya permission tertentu
// Otoritas Super Admin akan selalu lolos pengecekan ini
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, _ := c.Get("roles")
		userRoles := roles.([]string)

		permissions, _ := c.Get("permissions")
		userPermissions := permissions.([]string)

		// 1. Cek jika user adalah Super Admin (Bypass semuanya)
		for _, role := range userRoles {
			if role == "Super Admin" {
				c.Next()
				return
			}
		}

		// 2. Cek apakah user punya permission spesifik
		for _, p := range userPermissions {
			if p == permission {
				c.Next()
				return
			}
		}

		// Jika tidak punya akses
		response.Forbidden(c, "Anda tidak memiliki akses ke fitur ini ("+permission+")")
		c.Abort()
	}
}
