package sync

import (
	"github.com/alex/ads_backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.POST("/meta/sync", middleware.AuthMiddleware(), middleware.RequirePermission("meta.sync.trigger"), h.TriggerSync)
	r.GET("/meta/sync/status", middleware.AuthMiddleware(), middleware.RequirePermission("meta.sync.view"), h.SyncStatus)
	r.GET("/meta/sync/batches", middleware.AuthMiddleware(), middleware.RequirePermission("meta.sync.view"), h.ListBatches)
	r.GET("/meta/sync/batches/:id", middleware.AuthMiddleware(), middleware.RequirePermission("meta.sync.view"), h.GetBatch)
}
