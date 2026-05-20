package insight

import (
	"net/http"
	"strconv"

	"github.com/alex/ads_backend/internal/meta/insight/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

var _ = dto.InsightResponse{}

// GetCampaignInsights godoc
// @Summary      Get Campaign Insights
// @Description  Retrieve campaign-level insights from local database (synced from Meta)
// @Tags         Meta Insights
// @Accept       json
// @Produce      json
// @Param        account_id     query     string  false  "Filter by account ID"
// @Param        campaign_id    query     string  false  "Filter by campaign ID"
// @Param        date_start     query     string  false  "Filter by date start (YYYY-MM-DD)"
// @Param        date_stop      query     string  false  "Filter by date stop (YYYY-MM-DD)"
// @Param        page           query     int     false  "Page number" default(1)
// @Param        limit          query     int     false  "Items per page" default(25)
// @Success      200            {object}  response.Response{data=[]dto.InsightResponse,meta=response.PaginationMeta}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Router       /meta/insights/campaign [get]
func (h *Handler) GetCampaignInsights(c *gin.Context) {
	filter := InsightFilter{
		AccountID:  c.Query("account_id"),
		CampaignID: c.Query("campaign_id"),
		DateStart:  c.Query("date_start"),
		DateStop:   c.Query("date_stop"),
		Page:       parseQueryInt(c, "page", 1),
		Limit:      parseQueryInt(c, "limit", 25),
	}

	resp, meta, err := h.service.GetCampaignInsights(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.SuccessWithPagination(c, "Successfully retrieved campaign insights", resp, meta)
}

// GetAdInsights godoc
// @Summary      Get Ad Insights
// @Description  Retrieve ad-level insights from local database (synced from Meta)
// @Tags         Meta Insights
// @Accept       json
// @Produce      json
// @Param        account_id     query     string  false  "Filter by account ID"
// @Param        campaign_id    query     string  false  "Filter by campaign ID"
// @Param        adset_id       query     string  false  "Filter by adset ID"
// @Param        ad_id          query     string  false  "Filter by ad ID"
// @Param        date_start     query     string  false  "Filter by date start (YYYY-MM-DD)"
// @Param        date_stop      query     string  false  "Filter by date stop (YYYY-MM-DD)"
// @Param        page           query     int     false  "Page number" default(1)
// @Param        limit          query     int     false  "Items per page" default(25)
// @Success      200            {object}  response.Response{data=[]dto.InsightResponse,meta=response.PaginationMeta}
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Router       /meta/insights/ad [get]
func (h *Handler) GetAdInsights(c *gin.Context) {
	filter := InsightFilter{
		AccountID:  c.Query("account_id"),
		CampaignID: c.Query("campaign_id"),
		AdSetID:    c.Query("adset_id"),
		AdID:       c.Query("ad_id"),
		DateStart:  c.Query("date_start"),
		DateStop:   c.Query("date_stop"),
		Page:       parseQueryInt(c, "page", 1),
		Limit:      parseQueryInt(c, "limit", 25),
	}

	resp, meta, err := h.service.GetAdInsights(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.SuccessWithPagination(c, "Successfully retrieved ad insights", resp, meta)
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
