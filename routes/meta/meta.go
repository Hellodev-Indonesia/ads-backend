package meta

import (
	"github.com/alex/ads_backend/internal/meta"
	"github.com/gin-gonic/gin"
)

func RegisterMetaRoutes(r *gin.RouterGroup, h *meta.Handler) {
	meta.RegisterRoutes(r, h)
}
