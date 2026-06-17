package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/alex/ads_backend/internal/reporting/dashboard/dto"
	"github.com/alex/ads_backend/pkg/response"
)

type DashboardHandler struct {
	service DashboardService
}

var _ = dto.SummaryResponse{}
var _ = dto.BrandMonitoringResponse{}
var _ = dto.RecentActivityResponse{}

func NewDashboardHandler(service DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// GetSummary godoc
// @Summary      Get dashboard summary metrics
// @Description  Get total spend, ongoing ads, and security status
// @Tags         Reporting Dashboard
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Router       /reporting/dashboard/summary [get]
// @Security     BearerAuth
func (h *DashboardHandler) GetSummary(c *gin.Context) {
	ctx := c.Request.Context()
	
	res, err := h.service.GetSummary(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get dashboard summary", err)
		return
	}

	response.Success(c, "Success", res)
}

// GetBrandsMonitoring godoc
// @Summary      Get brands monitoring data
// @Description  Get metrics per brand including ad accounts, active campaigns, and total spends
// @Tags         Reporting Dashboard
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Router       /reporting/dashboard/brands [get]
// @Security     BearerAuth
func (h *DashboardHandler) GetBrandsMonitoring(c *gin.Context) {
	ctx := c.Request.Context()

	res, err := h.service.GetBrandsMonitoring(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get brands monitoring", err)
		return
	}

	response.Success(c, "Success", res)
}

// GetRecentActivities godoc
// @Summary      Get recent activities
// @Description  Get list of recent activities across brands
// @Tags         Reporting Dashboard
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Router       /reporting/dashboard/activities [get]
// @Security     BearerAuth
func (h *DashboardHandler) GetRecentActivities(c *gin.Context) {
	ctx := c.Request.Context()

	res, err := h.service.GetRecentActivities(ctx)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get recent activities", err)
		return
	}

	response.Success(c, "Success", res)
}
