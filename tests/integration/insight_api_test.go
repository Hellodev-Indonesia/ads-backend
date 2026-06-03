package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex/ads_backend/internal/meta/insight"
	"github.com/alex/ads_backend/internal/meta/insight/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupInsightAPI(t *testing.T) (*gin.Engine, *insight.MockService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := insight.NewMockService(t)
	handler := insight.NewHandler(mockService)

	insightGroup := router.Group("/meta/insights")
	{
		insightGroup.GET("/campaign", handler.GetCampaignInsights)
		insightGroup.GET("/ad", handler.GetAdInsights)
	}

	return router, mockService
}

func TestInsightAPI_GetCampaignInsights(t *testing.T) {
	router, mockService := setupInsightAPI(t)

	insights := []dto.InsightResponse{
		{CampaignID: "cmp_1", Spend: "100.00"},
	}
	meta := &response.PaginationMeta{Page: 1, Limit: 25, Total: 1, LastPage: 1}

	mockService.On("GetCampaignInsights", mock.MatchedBy(func(f insight.InsightFilter) bool {
		return f.Page == 1 && f.Limit == 25
	})).Return(insights, meta, nil)

	req, _ := http.NewRequest(http.MethodGet, "/meta/insights/campaign?page=1&limit=25", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInsightAPI_GetAdInsights(t *testing.T) {
	router, mockService := setupInsightAPI(t)

	insights := []dto.InsightResponse{
		{AdID: "ad_1", Spend: "50.00"},
	}
	meta := &response.PaginationMeta{Page: 1, Limit: 25, Total: 1, LastPage: 1}

	mockService.On("GetAdInsights", mock.MatchedBy(func(f insight.InsightFilter) bool {
		return f.Page == 1 && f.Limit == 25
	})).Return(insights, meta, nil)

	req, _ := http.NewRequest(http.MethodGet, "/meta/insights/ad?page=1&limit=25", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
