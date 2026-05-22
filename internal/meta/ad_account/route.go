package ad_account

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	group := r.Group("/meta/ad-accounts")
	group.Use(middleware.AuthMiddleware(), middleware.RequirePermission("meta.campaign.view"))
	{
		group.GET("", h.GetAdAccounts)
		group.GET("/unassigned", h.GetUnassignedAdAccounts)
		group.PUT("/:id/assign-brand", h.AssignBrand)
		group.PUT("/:id/unassign-brand", h.UnassignBrand)
	}
}
