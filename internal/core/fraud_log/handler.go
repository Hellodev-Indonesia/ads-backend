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

func (h *Handler) Resolve(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	log, err := h.service.Resolve(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, "Fraud log resolved", log)
}
