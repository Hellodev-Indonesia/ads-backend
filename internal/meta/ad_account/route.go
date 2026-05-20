package ad_account

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.GET("/meta/ad-accounts", middleware.AuthMiddleware(), middleware.RequirePermission("meta.campaign.view"), h.GetAdAccounts)
}
