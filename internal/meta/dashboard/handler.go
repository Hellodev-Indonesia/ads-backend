package dashboard

import (
	"net/http"
	"strconv"

	"github.com/alex/ads_backend/config"
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
// @Param        ad_account_id  query     string  false  "Ad Account ID (falls back to config)"
// @Param        status         query     string  false  "Filter by campaign status (ACTIVE, PAUSED, etc)"
// @Param        search         query     string  false  "Search by campaign name"
// @Param        page           query     int     false  "Page number" default(1)
// @Param        limit          query     int     false  "Items per page" default(25)
// @Success      200            {object}  response.Response
// @Failure      400            {object}  response.ErrorResponse
// @Failure      500            {object}  response.ErrorResponse
// @Router       /meta/dashboard/campaigns [get]
func (h *Handler) GetCampaignDashboard(c *gin.Context) {
	adAccountID := c.Query("ad_account_id")
	if adAccountID == "" {
		adAccountID = config.MetaAdAccountID
	}
	if adAccountID == "" {
		response.Error(c, http.StatusBadRequest, "Ad Account ID is required", nil)
		return
	}

	filter := DashboardFilter{
		AccountID: adAccountID,
		Status:    c.Query("status"),
		Search:    c.Query("search"),
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
