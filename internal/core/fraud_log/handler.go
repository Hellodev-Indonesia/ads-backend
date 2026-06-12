package fraud_log

import (
	"net/http"
	"strconv"

	"github.com/alex/ads_backend/internal/core/fraud_log/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// FindAll godoc
// @Summary      Get fraud logs
// @Description  Get a list of fraud logs with pagination and filters
// @Tags         Fraud Logs
// @Accept       json
// @Produce      json
// @Param        page          query     int     false  "Page number"
// @Param        limit         query     int     false  "Limit per page"
// @Param        brand_id      query     int     false  "Brand ID"
// @Param        ad_account_id query     string  false  "Ad Account ID"
// @Param        campaign_id   query     string  false  "Campaign ID"
// @Param        creative_id   query     string  false  "Creative ID"
// @Param        severity      query     string  false  "Severity"
// @Param        status        query     string  false  "Status"
// @Param        date_start    query     string  false  "Date Start"
// @Param        date_stop     query     string  false  "Date Stop"
// @Success      200  {object}  response.PaginationResponse{data=[]dto.FraudLogResponse}
// @Router       /fraud-logs [get]
// @Security     BearerAuth
func (h *Handler) FindAll(c *gin.Context) {
	var filter dto.FraudLogFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	logs, total, err := h.service.FindAll(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 25
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	lastPage := int(total) / limit
	if int(total)%limit > 0 {
		lastPage++
	}

	response.SuccessWithPagination(c, "Successfully retrieved fraud logs", logs, &response.PaginationMeta{
		Page:     page,
		Limit:    limit,
		Total:    int(total),
		LastPage: lastPage,
	})
}

// FindByID godoc
// @Summary      Get fraud log by ID
// @Description  Get a fraud log by its ID
// @Tags         Fraud Logs
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Fraud Log ID"
// @Success      200  {object}  response.SuccessResponse{data=dto.FraudLogResponse}
// @Router       /fraud-logs/{id} [get]
// @Security     BearerAuth
func (h *Handler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	log, err := h.service.FindByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, "Successfully retrieved fraud log", log)
}

// Resolve godoc
// @Summary      Resolve a fraud log
// @Description  Mark a fraud log as resolved
// @Tags         Fraud Logs
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Fraud Log ID"
// @Success      200  {object}  response.SuccessResponse{data=dto.FraudLogResponse}
// @Router       /fraud-logs/{id}/resolve [put]
// @Security     BearerAuth
func (h *Handler) Resolve(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	userID := uint64(userIDVal.(float64))

	log, err := h.service.Resolve(id, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, "Fraud log resolved", log)
}
