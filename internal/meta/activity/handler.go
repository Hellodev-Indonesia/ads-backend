package activity

import (
	"net/http"
	"strconv"

	"github.com/alex/ads_backend/internal/meta/activity/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// GetAllActivities godoc
// @Summary      Get All Activities
// @Description  Retrieve meta activities across all ad accounts
// @Tags         Meta Activities
// @Accept       json
// @Produce      json
// @Param        page           query     int     false  "Page number" default(1)
// @Param        limit          query     int     false  "Items per page" default(25)
// @Success      200            {object}  response.Response{data=[]dto.ActivityResponse,meta=response.PaginationMeta}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Security BearerAuth
// @Router       /meta/activities [get]
func (h *Handler) GetAllActivities(c *gin.Context) {
	var filter dto.ActivityFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	resp, total, err := h.service.FindAll(filter)
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

	meta := &response.PaginationMeta{
		Page:     page,
		Limit:    limit,
		Total:    int(total),
		LastPage: lastPage,
	}

	response.SuccessWithPagination(c, "Successfully retrieved all activities", resp, meta)
}

// GetActivitiesByBrand godoc
// @Summary      Get Activities by Brand ID
// @Description  Retrieve meta activities for all ad accounts under a specific brand
// @Tags         Meta Activities
// @Accept       json
// @Produce      json
// @Param        brand_id       path      int     true   "Brand ID"
// @Param        page           query     int     false  "Page number" default(1)
// @Param        limit          query     int     false  "Items per page" default(25)
// @Success      200            {object}  response.Response{data=[]dto.ActivityResponse,meta=response.PaginationMeta}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Security BearerAuth
// @Router       /meta/brands/{brand_id}/activities [get]
func (h *Handler) GetActivitiesByBrand(c *gin.Context) {
	brandIDParam := c.Param("brand_id")
	brandID, err := strconv.ParseUint(brandIDParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid brand ID", nil)
		return
	}

	var filter dto.ActivityFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	resp, total, err := h.service.FindAllByBrand(brandID, filter)
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

	meta := &response.PaginationMeta{
		Page:     page,
		Limit:    limit,
		Total:    int(total),
		LastPage: lastPage,
	}

	response.SuccessWithPagination(c, "Successfully retrieved activities", resp, meta)
}
