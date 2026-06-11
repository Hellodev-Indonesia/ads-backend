package sync

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/alex/ads_backend/internal/meta/sync/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// JobTrigger is the subset of MetaAdsSyncJob the handler needs.
type JobTrigger interface {
	Start(ctx context.Context, req dto.TriggerSyncRequest) ([]*MetaSyncBatch, error)
	IsRunning() bool
}

type Handler struct {
	job     JobTrigger
	service *Service
}

func NewHandler(job JobTrigger, service *Service) *Handler {
	return &Handler{job: job, service: service}
}

// TriggerSync godoc
// @Summary      Trigger Meta Ads Sync
// @Description  Manually trigger a full Meta Ads sync. Returns immediately; subscribe to the Centrifugo channel for real-time progress.
// @Tags         Meta Sync
// @Param        request  body      dto.TriggerSyncRequest  false  "Sync Request (optional ad_account_id, date_start, date_stop)"
// @Success      202  {object}  response.Response
// @Failure      409  {object}  response.ErrorResponse  "Sync already in progress"
// @Failure      500  {object}  response.ErrorResponse
// @Security     BearerAuth
// @Security BearerAuth
// @Router       /meta/sync [post]
func (h *Handler) TriggerSync(c *gin.Context) {
	var req dto.TriggerSyncRequest
	if err := c.ShouldBind(&req); err != nil {
		// Log but continue, allowing fallback to query parameters or empty
		_ = c.ShouldBindQuery(&req)
	}
	if req.AdAccountID == "" {
		req.AdAccountID = c.Query("ad_account_id")
	}

	batches, err := h.job.Start(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrAlreadyRunning) {
			response.Error(c, http.StatusConflict, "Sync is already in progress", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var batchIDs []uint64
	var batchCodes []string
	for _, b := range batches {
		batchIDs = append(batchIDs, b.ID)
		batchCodes = append(batchCodes, b.BatchCode)
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Sync started",
		"data": gin.H{
			"batch_ids":   batchIDs,
			"batch_codes": batchCodes,
			"channel":     Channel,
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
// @Security BearerAuth
// @Router       /meta/sync/status [get]
func (h *Handler) SyncStatus(c *gin.Context) {
	respData := gin.H{
		"running": h.job.IsRunning(),
		"channel": Channel,
	}

	lastBatch, err := h.service.GetLastSyncBatch(c.Request.Context())
	if err == nil && lastBatch != nil && lastBatch.FinishedAt != nil {
		minutesAgo := int(time.Since(*lastBatch.FinishedAt).Minutes())
		respData["last_sync_minutes_ago"] = minutesAgo
		respData["last_sync_at"] = lastBatch.FinishedAt
	}

	response.Success(c, "OK", respData)
}

// ListBatches godoc
// @Summary      List Sync Batches
// @Description  Returns a paginated list of sync batches ordered by most recent first.
// @Tags         Meta Sync
// @Produce      json
// @Param        page   query     int  false  "Page number" default(1)
// @Param        limit  query     int  false  "Items per page" default(25)
// @Success      200    {object}  response.PaginationResponse{data=[]MetaSyncBatch}
// @Failure      500    {object}  response.ErrorResponse
// @Security     BearerAuth
// @Security BearerAuth
// @Router       /meta/sync/batches [get]
func (h *Handler) ListBatches(c *gin.Context) {
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 25)

	batches, total, err := h.service.ListBatches(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	lastPage := int(total) / limit
	if int(total)%limit > 0 {
		lastPage++
	}

	response.SuccessWithPagination(c, "Successfully retrieved sync batches", batches, &response.PaginationMeta{
		Page:     page,
		Limit:    limit,
		Total:    int(total),
		LastPage: lastPage,
	})
}

// GetBatch godoc
// @Summary      Get Sync Batch
// @Description  Returns a single sync batch with its steps.
// @Tags         Meta Sync
// @Produce      json
// @Param        id   path      int  true  "Batch ID"
// @Success      200  {object}  response.SuccessResponse{data=MetaSyncBatch}
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Security     BearerAuth
// @Security BearerAuth
// @Router       /meta/sync/batches/{id} [get]
func (h *Handler) GetBatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid batch ID", nil)
		return
	}

	batch, err := h.service.GetBatchByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "batch not found", nil)
		return
	}

	response.Success(c, "Successfully retrieved sync batch", batch)
}

func parseIntQuery(c *gin.Context, key string, defaultVal int) int {
	v, err := strconv.Atoi(c.Query(key))
	if err != nil || v <= 0 {
		return defaultVal
	}
	return v
}
