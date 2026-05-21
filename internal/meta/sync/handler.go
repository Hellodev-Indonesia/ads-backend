package sync

import (
	"context"
	"errors"
	"net/http"

	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// JobTrigger is the subset of MetaAdsSyncJob the handler needs.
type JobTrigger interface {
	Start(ctx context.Context) (*MetaSyncBatch, error)
	IsRunning() bool
}

type Handler struct {
	job JobTrigger
}

func NewHandler(job JobTrigger) *Handler {
	return &Handler{job: job}
}

// TriggerSync godoc
// @Summary      Trigger Meta Ads Sync
// @Description  Manually trigger a full Meta Ads sync. Returns immediately; subscribe to the Centrifugo channel for real-time progress.
// @Tags         Meta Sync
// @Produce      json
// @Success      202  {object}  response.Response
// @Failure      409  {object}  response.ErrorResponse  "Sync already in progress"
// @Failure      500  {object}  response.ErrorResponse
// @Security     BearerAuth
// @Router       /meta/sync [post]
func (h *Handler) TriggerSync(c *gin.Context) {
	batch, err := h.job.Start(c.Request.Context())
	if err != nil {
		if errors.Is(err, ErrAlreadyRunning) {
			response.Error(c, http.StatusConflict, "Sync is already in progress", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Sync started",
		"data": gin.H{
			"batch_id":   batch.ID,
			"batch_code": batch.BatchCode,
			"channel":    Channel,
		},
	})
}

// SyncStatus godoc
// @Summary      Get Sync Status
// @Description  Returns whether a sync job is currently running.
// @Tags         Meta Sync
// @Produce      json
// @Success      200  {object}  response.Response
// @Security     BearerAuth
// @Router       /meta/sync/status [get]
func (h *Handler) SyncStatus(c *gin.Context) {
	response.Success(c, "OK", gin.H{
		"running": h.job.IsRunning(),
		"channel": Channel,
	})
}
