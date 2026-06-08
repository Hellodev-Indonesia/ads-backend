package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex/ads_backend/internal/meta/campaign"
	"github.com/alex/ads_backend/internal/meta/campaign/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupCampaignAPI(t *testing.T) (*gin.Engine, *campaign.MockService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := campaign.NewMockService(t)
	handler := campaign.NewHandler(mockService)

	campaignGroup := router.Group("/meta/campaigns")
	{
		campaignGroup.GET("", handler.GetCampaigns)
	}

	return router, mockService
}

func TestCampaignAPI_GetCampaigns(t *testing.T) {
	router, mockService := setupCampaignAPI(t)

	campaigns := []dto.CampaignResponse{
		{ID: "cmp_1", Name: "Test Campaign", Status: "ACTIVE"},
	}
	meta := &response.PaginationMeta{Page: 1, Limit: 25, Total: 1, LastPage: 1}

	mockService.On("GetCampaigns", mock.MatchedBy(func(f campaign.CampaignFilter) bool {
		return f.Page == 1 && f.Limit == 25
	})).Return(campaigns, meta, nil)

	req, _ := http.NewRequest(http.MethodGet, "/meta/campaigns?page=1&limit=25", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
