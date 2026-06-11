package ads

import (
	"net/http"
	"strconv"

	"github.com/alex/ads_backend/internal/meta/ads/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

var _ = dto.AdResponse{}
var _ = dto.CreativeResponse{}

// GetAds godoc
// @Summary      Get Ads
// @Description  Retrieve ads from local database (synced from Meta)
// @Tags         Meta Ads
// @Accept       json
// @Produce      json
// @Param        campaign_id    query     string  false  "Filter by campaign ID"
// @Param        adset_id       query     string  false  "Filter by adset ID"
// @Param        status         query     string  false  "Filter by status (ACTIVE, PAUSED, etc)"
// @Param        search         query     string  false  "Search by ad name"
// @Param        page           query     int     false  "Page number" default(1)
// @Param        limit          query     int     false  "Items per page" default(25)
// @Success      200            {object}  response.Response{data=[]dto.AdResponse,meta=response.PaginationMeta}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Security BearerAuth
// @Router       /meta/ads [get]
func (h *Handler) GetAds(c *gin.Context) {
	filter := AdFilter{
		CampaignID: c.Query("campaign_id"),
		AdSetID:    c.Query("adset_id"),
		Status:     c.Query("status"),
		Search:     c.Query("search"),
		Page:       parseQueryInt(c, "page", 1),
		Limit:      parseQueryInt(c, "limit", 25),
	}

	resp, meta, err := h.service.GetAds(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.SuccessWithPagination(c, "Successfully retrieved ads", resp, meta)
}

// GetAdsByBrand godoc
// @Summary      Get Ads by Brand ID
// @Description  Retrieve ads for a specific brand
// @Tags         Meta Ads
// @Accept       json
// @Produce      json
// @Param        brand_id       path      int     true   "Brand ID"
// @Param        campaign_id    query     string  false  "Filter by campaign ID"
// @Param        adset_id       query     string  false  "Filter by adset ID"
// @Param        status         query     string  false  "Filter by status (ACTIVE, PAUSED, etc)"
// @Param        search         query     string  false  "Search by ad name"
// @Param        page           query     int     false  "Page number" default(1)
// @Param        limit          query     int     false  "Items per page" default(25)
// @Param        date_start     query     string  false  "Start Date (YYYY-MM-DD)"
// @Param        date_stop      query     string  false  "Stop Date (YYYY-MM-DD)"
// @Success      200            {object}  response.Response{data=[]dto.AdDashboardRow,meta=response.PaginationMeta}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Security BearerAuth
// @Router       /meta/brands/{brand_id}/ads [get]
func (h *Handler) GetAdsByBrand(c *gin.Context) {
	brandIDParam := c.Param("brand_id")
	brandID, err := strconv.ParseUint(brandIDParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid brand ID", nil)
		return
	}

	filter := AdFilter{
		BrandID:     &brandID,
		CampaignID:  c.Query("campaign_id"),
		AdSetID:     c.Query("adset_id"),
		CampaignIDs: c.QueryArray("campaign_ids"),
		AdSetIDs:    c.QueryArray("adset_ids"),
		Status:      c.Query("status"),
		Search:      c.Query("search"),
		DateStart:   c.Query("date_start"),
		DateStop:    c.Query("date_stop"),
		Page:        parseQueryInt(c, "page", 1),
		Limit:       parseQueryInt(c, "limit", 25),
	}

	resp, meta, err := h.service.GetAdDashboard(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.SuccessWithPagination(c, "Successfully retrieved ads for brand", resp, meta)
}

// GetCreative godoc
// @Summary      Get Ad Creative
// @Description  Retrieve details of a specific ad creative from the local database (synced from Meta)
// @Tags         Meta Ads
// @Accept       json
// @Produce      json
// @Param        id       path      string  true   "Creative ID"
// @Param        fields   query     string  false  "Custom fields comma-separated preset"
// @Success      200      {object}  response.Response{data=dto.CreativeResponse}
// @Failure      400      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Security BearerAuth
// @Router       /meta/creatives/{id} [get]
func (h *Handler) GetCreative(c *gin.Context) {
	creativeID := c.Param("id")
	if creativeID == "" {
		response.Error(c, http.StatusBadRequest, "Creative ID is required", nil)
		return
	}

	fields := c.Query("fields")
	if fields == "" {
		fields = DefaultCreativeFields
	}

	resp, err := h.service.GetCreative(creativeID, fields)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, "Successfully retrieved ad creative", resp)
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

// GetAdDashboard godoc
// @Summary      Ad Dashboard
// @Description  Returns ads with performance metrics joined from insights
// @Tags         Meta Ads
// @Accept       json
// @Produce      json
// @Param        ad_account_id  query     string  false  "Ad Account ID"
// @Param        brand_id       query     int     false  "Brand ID"
// @Param        campaign_id    query     string  false  "Filter by campaign ID"
// @Param        adset_id       query     string  false  "Filter by ad set ID"
// @Param        status         query     string  false  "Filter by ad status"
// @Param        search         query     string  false  "Search by ad name"
// @Param        date_start     query     string  false  "Filter by date start (YYYY-MM-DD)"
// @Param        date_stop      query     string  false  "Filter by date stop (YYYY-MM-DD)"
// @Param        page           query     int     false  "Page number" default(1)
// @Param        limit          query     int     false  "Items per page" default(25)
// @Success      200            {object}  response.Response
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Security BearerAuth
// @Router       /meta/ads/dashboard [get]
func (h *Handler) GetAdDashboard(c *gin.Context) {
	var brandID *uint64
	if brandIDStr := c.Query("brand_id"); brandIDStr != "" {
		bid, err := strconv.ParseUint(brandIDStr, 10, 64)
		if err == nil {
			brandID = &bid
		}
	}

	filter := AdFilter{
		AccountID:   c.Query("ad_account_id"),
		BrandID:     brandID,
		CampaignID:  c.Query("campaign_id"),
		AdSetID:     c.Query("adset_id"),
		CampaignIDs: c.QueryArray("campaign_ids"),
		AdSetIDs:    c.QueryArray("adset_ids"),
		Status:      c.Query("status"),
		Search:      c.Query("search"),
		DateStart:   c.Query("date_start"),
		DateStop:    c.Query("date_stop"),
		Page:        parseQueryInt(c, "page", 1),
		Limit:       parseQueryInt(c, "limit", 25),
	}

	rows, meta, err := h.service.GetAdDashboard(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.SuccessWithPagination(c, "Successfully retrieved ad dashboard", rows, meta)
}
