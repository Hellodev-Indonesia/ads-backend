package dashboard

import (
	"net/http"
	"strconv"

	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// GetCampaignDashboard godoc
// @Summary      Campaign Dashboard
// @Description  Returns campaigns with performance metrics (spend, impressions, reach, actions) joined from insights and adsets
// @Tags         Meta Dashboard
// @Accept       json
// @Produce      json
// @Param        ad_account_id  query     string  false  "Ad Account ID"
// @Param        brand_id       query     int     false  "Brand ID"
// @Param        status         query     string  false  "Filter by campaign status (ACTIVE, PAUSED, etc)"
// @Param        search         query     string  false  "Search by campaign name"
// @Param        date_start     query     string  false  "Filter by date start (YYYY-MM-DD)"
// @Param        date_stop      query     string  false  "Filter by date stop (YYYY-MM-DD)"
// @Param        page           query     int     false  "Page number" default(1)
// @Param        limit          query     int     false  "Items per page" default(25)
// @Success      200            {object}  response.Response
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Router       /meta/dashboard/campaigns [get]
func (h *Handler) GetCampaignDashboard(c *gin.Context) {
	adAccountID := c.Query("ad_account_id")

	var brandID *uint64
	brandIDStr := c.Query("brand_id")
	if brandIDStr != "" {
		bid, err := strconv.ParseUint(brandIDStr, 10, 64)
		if err == nil {
			brandID = &bid
		}
	}

	filter := DashboardFilter{
		AccountID: adAccountID,
		BrandID:   brandID,
		Status:    c.Query("status"),
		Search:    c.Query("search"),
		DateStart: c.Query("date_start"),
		DateStop:  c.Query("date_stop"),
		Page:      parseQueryInt(c, "page", 1),
		Limit:     parseQueryInt(c, "limit", 25),
	}

	rows, meta, err := h.service.GetCampaignDashboard(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.SuccessWithPagination(c, "Successfully retrieved campaign dashboard", rows, meta)
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
// @Tags         Meta Dashboard
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
// @Router       /meta/dashboard/adsets [get]
func (h *Handler) GetAdSetDashboard(c *gin.Context) {
	var brandID *uint64
	if brandIDStr := c.Query("brand_id"); brandIDStr != "" {
		bid, err := strconv.ParseUint(brandIDStr, 10, 64)
		if err == nil {
			brandID = &bid
		}
	}

	filter := DashboardFilter{
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

// GetAdDashboard godoc
// @Summary      Ad Dashboard
// @Description  Returns ads with performance metrics joined from insights
// @Tags         Meta Dashboard
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
// @Router       /meta/dashboard/ads [get]
func (h *Handler) GetAdDashboard(c *gin.Context) {
	var brandID *uint64
	if brandIDStr := c.Query("brand_id"); brandIDStr != "" {
		bid, err := strconv.ParseUint(brandIDStr, 10, 64)
		if err == nil {
			brandID = &bid
		}
	}

	filter := DashboardFilter{
		AccountID:  c.Query("ad_account_id"),
		BrandID:    brandID,
		CampaignID: c.Query("campaign_id"),
		AdSetID:    c.Query("adset_id"),
		Status:     c.Query("status"),
		Search:     c.Query("search"),
		DateStart:  c.Query("date_start"),
		DateStop:   c.Query("date_stop"),
		Page:       parseQueryInt(c, "page", 1),
		Limit:      parseQueryInt(c, "limit", 25),
	}

	rows, meta, err := h.service.GetAdDashboard(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.SuccessWithPagination(c, "Successfully retrieved ad dashboard", rows, meta)
}
