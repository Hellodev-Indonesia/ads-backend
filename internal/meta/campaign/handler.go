package campaign

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/alex/ads_backend/internal/meta/campaign/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

var _ = dto.CampaignResponse{}

// GetCampaigns godoc
// @Summary      Get Campaigns
// @Description  Retrieve campaigns from local database (synced from Meta)
// @Tags         Meta Campaigns
// @Accept       json
// @Produce      json
// @Param        ad_account_id  query     string  false  "Ad Account ID (falls back to config.MetaAdAccountID)"
// @Param        status         query     string  false  "Filter by status (ACTIVE, PAUSED, etc)"
// @Param        search         query     string  false  "Search by campaign name"
// @Param        page           query     int     false  "Page number" default(1)
// @Param        limit          query     int     false  "Items per page" default(25)
// @Success      200            {object}  response.Response{data=[]dto.CampaignResponse,meta=response.PaginationMeta}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Security BearerAuth
// @Router       /meta/campaigns [get]
func (h *Handler) GetCampaigns(c *gin.Context) {
	adAccountID := c.Query("ad_account_id")

	filter := CampaignFilter{
		AccountID: adAccountID,
		Status:    c.Query("status"),
		Search:    c.Query("search"),
		Page:      parseQueryInt(c, "page", 1),
		Limit:     parseQueryInt(c, "limit", 25),
	}

	resp, meta, err := h.service.GetCampaigns(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.SuccessWithPagination(c, "Successfully retrieved campaigns", resp, meta)
}

// GetCampaignsByBrand godoc
// @Summary      Get Campaigns by Brand ID
// @Description  Retrieve campaigns for a specific brand
// @Tags         Meta Campaigns
// @Accept       json
// @Produce      json
// @Param        brand_id       path      int     true   "Brand ID"
// @Param        status         query     string  false  "Filter by status (ACTIVE, PAUSED, etc)"
// @Param        search         query     string  false  "Search by campaign name"
// @Param        page           query     int     false  "Page number" default(1)
// @Param        limit          query     int     false  "Items per page" default(25)
// @Param        date_start     query     string  false  "Start Date (YYYY-MM-DD)"
// @Param        date_stop      query     string  false  "Stop Date (YYYY-MM-DD)"
// @Success      200            {object}  response.Response{data=[]dto.CampaignDashboardRow,meta=response.PaginationMeta}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Security BearerAuth
// @Router       /meta/brands/{brand_id}/campaigns [get]
func (h *Handler) GetCampaignsByBrand(c *gin.Context) {
	brandIDParam := c.Param("brand_id")
	brandID, err := strconv.ParseUint(brandIDParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid brand ID", nil)
		return
	}

	filter := CampaignFilter{
		BrandID:     &brandID,
		Status:      c.Query("status"),
		Search:      c.Query("search"),
		DateStart:   c.Query("date_start"),
		DateStop:    c.Query("date_stop"),
		CampaignIDs: parseQueryArray(c, "campaign_ids"),
		AdSetIDs:    parseQueryArray(c, "adset_ids"),
		Page:        parseQueryInt(c, "page", 1),
		Limit:       parseQueryInt(c, "limit", 25),
	}

	resp, meta, err := h.service.GetCampaignDashboard(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.SuccessWithPagination(c, "Successfully retrieved campaigns for brand", resp, meta)
}

// GetCampaignSummaryByBrand godoc
// @Summary      Get Campaign Summary by Brand ID
// @Description  Retrieve summary metrics for all campaigns under a specific brand
// @Tags         Meta Campaigns
// @Accept       json
// @Produce      json
// @Param        brand_id       path      int     true   "Brand ID"
// @Param        date_start     query     string  false  "Start Date (YYYY-MM-DD)"
// @Param        date_stop      query     string  false  "Stop Date (YYYY-MM-DD)"
// @Param        campaign_ids   query     []string false "Filter by Campaign IDs"
// @Param        adset_ids      query     []string false "Filter by AdSet IDs"
// @Success      200            {object}  response.Response{data=dto.CampaignSummaryResponse}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Security BearerAuth
// @Router       /meta/brands/{brand_id}/campaigns/summary [get]
func (h *Handler) GetCampaignSummaryByBrand(c *gin.Context) {
	brandIDParam := c.Param("brand_id")
	brandID, err := strconv.ParseUint(brandIDParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid brand ID", nil)
		return
	}

	dateStart := c.Query("date_start")
	dateStop := c.Query("date_stop")
	campaignIDs := c.QueryArray("campaign_ids")
	adsetIDs := c.QueryArray("adset_ids")

	resp, err := h.service.GetSummaryByBrand(brandID, dateStart, dateStop, campaignIDs, adsetIDs)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(c, "Successfully retrieved campaign summary for brand", resp)
}

// GetCampaignListByBrand godoc
// @Summary      Get Campaign Simple List by Brand ID
// @Description  Retrieve id and name of all campaigns under a specific brand
// @Tags         Meta Campaigns
// @Accept       json
// @Produce      json
// @Param        brand_id       path      int     true   "Brand ID"
// @Success      200            {object}  response.Response{data=[]dto.SimpleListResponse}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Security BearerAuth
// @Router       /meta/brands/{brand_id}/campaigns/list [get]
func (h *Handler) GetCampaignListByBrand(c *gin.Context) {
	brandIDParam := c.Param("brand_id")
	brandID, err := strconv.ParseUint(brandIDParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid brand ID", nil)
		return
	}

	resp, err := h.service.GetCampaignListByBrand(brandID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(c, "Successfully retrieved campaign list for brand", resp)
}

// GetCampaignDashboard godoc
// @Summary      Campaign Dashboard
// @Description  Returns campaigns with performance metrics (spend, impressions, reach, actions) joined from insights and adsets
// @Tags         Meta Campaigns
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
// @Success      200            {object}  response.Response{data=[]dto.CampaignDashboardRow,meta=response.PaginationMeta}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Security BearerAuth
// @Router       /meta/campaigns/dashboard [get]
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

	filter := CampaignFilter{
		AccountID:   adAccountID,
		BrandID:     brandID,
		Status:      c.Query("status"),
		Search:      c.Query("search"),
		DateStart:   c.Query("date_start"),
		DateStop:    c.Query("date_stop"),
		CampaignIDs: parseQueryArray(c, "campaign_ids"),
		AdSetIDs:    parseQueryArray(c, "adset_ids"),
		Page:        parseQueryInt(c, "page", 1),
		Limit:       parseQueryInt(c, "limit", 25),
	}

	resp, meta, err := h.service.GetCampaignDashboard(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.SuccessWithPagination(c, "Successfully retrieved campaign dashboard", resp, meta)
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

func parseQueryArray(c *gin.Context, key string) []string {
	vals := c.QueryArray(key)
	
	// If the array is empty, check if it's passed as a single string (query params without [])
	if len(vals) == 0 {
		val := c.Query(key)
		if val != "" {
			vals = []string{val}
		}
	}
	
	var result []string
	for _, v := range vals {
		if strings.Contains(v, ",") {
			parts := strings.Split(v, ",")
			for _, p := range parts {
				trimmed := strings.TrimSpace(p)
				if trimmed != "" {
					result = append(result, trimmed)
				}
			}
		} else {
			trimmed := strings.TrimSpace(v)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
	}
	return result
}
