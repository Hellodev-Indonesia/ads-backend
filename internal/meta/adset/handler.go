package adset

import (
	"net/http"
	"strconv"

	"github.com/alex/ads_backend/internal/meta/adset/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

var _ = dto.AdSetResponse{}

// GetAdSets godoc
// @Summary      Get AdSets
// @Description  Retrieve adsets from local database (synced from Meta)
// @Tags         Meta AdSets
// @Accept       json
// @Produce      json
// @Param        campaign_id    query     string  false  "Filter by campaign ID"
// @Param        status         query     string  false  "Filter by status (ACTIVE, PAUSED, etc)"
// @Param        search         query     string  false  "Search by adset name"
// @Param        page           query     int     false  "Page number" default(1)
// @Param        limit          query     int     false  "Items per page" default(25)
// @Success      200            {object}  response.Response{data=[]dto.AdSetResponse,meta=response.PaginationMeta}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Security BearerAuth
// @Router       /meta/adsets [get]
func (h *Handler) GetAdSets(c *gin.Context) {
	filter := AdSetFilter{
		CampaignID: c.Query("campaign_id"),
		Status:     c.Query("status"),
		Search:     c.Query("search"),
		Page:       parseQueryInt(c, "page", 1),
		Limit:      parseQueryInt(c, "limit", 25),
	}

	resp, meta, err := h.service.GetAdSets(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.SuccessWithPagination(c, "Successfully retrieved adsets", resp, meta)
}

// GetAdSetsByBrand godoc
// @Summary      Get AdSets by Brand ID
// @Description  Retrieve adsets for a specific brand
// @Tags         Meta AdSets
// @Accept       json
// @Produce      json
// @Param        brand_id       path      int     true   "Brand ID"
// @Param        campaign_id    query     string  false  "Filter by campaign ID"
// @Param        status         query     string  false  "Filter by status (ACTIVE, PAUSED, etc)"
// @Param        search         query     string  false  "Search by adset name"
// @Param        page           query     int     false  "Page number" default(1)
// @Param        limit          query     int     false  "Items per page" default(25)
// @Param        date_start     query     string  false  "Start Date (YYYY-MM-DD)"
// @Param        date_stop      query     string  false  "Stop Date (YYYY-MM-DD)"
// @Success      200            {object}  response.Response{data=[]dto.AdSetDashboardRow,meta=response.PaginationMeta}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Security BearerAuth
// @Router       /meta/brands/{brand_id}/adsets [get]
func (h *Handler) GetAdSetsByBrand(c *gin.Context) {
	brandIDParam := c.Param("brand_id")
	brandID, err := strconv.ParseUint(brandIDParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid brand ID", nil)
		return
	}

	filter := AdSetFilter{
		BrandID:    &brandID,
		CampaignID: c.Query("campaign_id"),
		Status:     c.Query("status"),
		Search:     c.Query("search"),
		DateStart:  c.Query("date_start"),
		DateStop:   c.Query("date_stop"),
		Page:       parseQueryInt(c, "page", 1),
		Limit:      parseQueryInt(c, "limit", 25),
	}

	resp, meta, err := h.service.GetAdSetDashboard(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.SuccessWithPagination(c, "Successfully retrieved adsets for brand", resp, meta)
}

func parseQueryInt(c *gin.Context, key string, defaultVal int) int {
	val := c.Query(key)
	if val == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(val)
	if err != nil || v <= 0 {
		return defaultVal
	}
	return v
}

// GetAdSetDashboard godoc
// @Summary      Ad Set Dashboard
// @Description  Returns ad sets with performance metrics joined from insights
// @Tags         Meta AdSets
// @Accept       json
// @Produce      json
// @Param        ad_account_id  query     string  false  "Ad Account ID"
// @Param        brand_id       query     int     false  "Brand ID"
// @Param        campaign_id    query     string  false  "Filter by campaign ID"
// @Param        status         query     string  false  "Filter by ad set status"
// @Param        search         query     string  false  "Search by ad set name"
// @Param        date_start     query     string  false  "Filter by date start (YYYY-MM-DD)"
// @Param        date_stop      query     string  false  "Filter by date stop (YYYY-MM-DD)"
// @Param        page           query     int     false  "Page number" default(1)
// @Param        limit          query     int     false  "Items per page" default(25)
// @Success      200            {object}  response.Response
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Security BearerAuth
// @Router       /meta/adsets/dashboard [get]
func (h *Handler) GetAdSetDashboard(c *gin.Context) {
	var brandID *uint64
	if brandIDStr := c.Query("brand_id"); brandIDStr != "" {
		bid, err := strconv.ParseUint(brandIDStr, 10, 64)
		if err == nil {
			brandID = &bid
		}
	}

	filter := AdSetFilter{
		AccountID:  c.Query("ad_account_id"),
		BrandID:    brandID,
		CampaignID: c.Query("campaign_id"),
		Status:     c.Query("status"),
		Search:     c.Query("search"),
		DateStart:  c.Query("date_start"),
		DateStop:   c.Query("date_stop"),
		Page:       parseQueryInt(c, "page", 1),
		Limit:      parseQueryInt(c, "limit", 25),
	}

	rows, meta, err := h.service.GetAdSetDashboard(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.SuccessWithPagination(c, "Successfully retrieved adset dashboard", rows, meta)
}
