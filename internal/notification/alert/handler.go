package alert

import (
	"net/http"
	"strconv"

	"github.com/alex/ads_backend/internal/notification/alert/dto"
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
	var filter dto.AlertFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	alerts, total, err := h.service.FindAll(filter)
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

	response.SuccessWithPagination(c, "Successfully retrieved alerts", alerts, &response.PaginationMeta{
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

	a, err := h.service.FindByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, "Successfully retrieved alert", a)
}

func (h *Handler) MarkRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	a, err := h.service.MarkRead(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, "Alert marked as read", a)
}
