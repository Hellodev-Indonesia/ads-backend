package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex/ads_backend/internal/meta/dashboard"
	"github.com/alex/ads_backend/internal/meta/dashboard/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupDashboardAPI(t *testing.T) (*gin.Engine, *dashboard.MockService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := dashboard.NewMockService(t)
	handler := dashboard.NewHandler(mockService)

	dashboardGroup := router.Group("/meta/dashboard")
	{
		dashboardGroup.GET("/campaigns", handler.GetCampaignDashboard)
		dashboardGroup.GET("/adsets", handler.GetAdSetDashboard)
		dashboardGroup.GET("/ads", handler.GetAdDashboard)
	}

	return router, mockService
}

func TestDashboardAPI_GetCampaignDashboard(t *testing.T) {
	router, mockService := setupDashboardAPI(t)

	rows := []dto.CampaignDashboardRow{
		{CampaignID: "cmp_1", CampaignName: "Campaign 1"},
	}
	meta := &response.PaginationMeta{Page: 1, Limit: 25, Total: 1, LastPage: 1}

	mockService.On("GetCampaignDashboard", mock.MatchedBy(func(f dashboard.DashboardFilter) bool {
		return f.Page == 1 && f.Limit == 25
	})).Return(rows, meta, nil)

	req, _ := http.NewRequest(http.MethodGet, "/meta/dashboard/campaigns?page=1&limit=25", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDashboardAPI_GetAdSetDashboard(t *testing.T) {
	router, mockService := setupDashboardAPI(t)

	rows := []dto.AdSetDashboardRow{
		{AdSetID: "adset_1", AdSetName: "AdSet 1"},
	}
	meta := &response.PaginationMeta{Page: 1, Limit: 25, Total: 1, LastPage: 1}

	mockService.On("GetAdSetDashboard", mock.MatchedBy(func(f dashboard.DashboardFilter) bool {
		return f.Page == 1 && f.Limit == 25
	})).Return(rows, meta, nil)

	req, _ := http.NewRequest(http.MethodGet, "/meta/dashboard/adsets?page=1&limit=25", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDashboardAPI_GetAdDashboard(t *testing.T) {
	router, mockService := setupDashboardAPI(t)

	rows := []dto.AdDashboardRow{
		{AdID: "ad_1", AdName: "Ad 1"},
	}
	meta := &response.PaginationMeta{Page: 1, Limit: 25, Total: 1, LastPage: 1}

	mockService.On("GetAdDashboard", mock.MatchedBy(func(f dashboard.DashboardFilter) bool {
		return f.Page == 1 && f.Limit == 25
	})).Return(rows, meta, nil)

	req, _ := http.NewRequest(http.MethodGet, "/meta/dashboard/ads?page=1&limit=25", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
