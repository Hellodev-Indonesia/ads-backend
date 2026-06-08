package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex/ads_backend/internal/meta/ads"
	"github.com/alex/ads_backend/internal/meta/ads/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAdsAPI(t *testing.T) (*gin.Engine, *ads.MockService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := ads.NewMockService(t)
	handler := ads.NewHandler(mockService)

	adsGroup := router.Group("/meta")
	{
		adsGroup.GET("/ads", handler.GetAds)
		adsGroup.GET("/creatives/:id", handler.GetCreative)
	}

	return router, mockService
}

func TestAdsAPI_GetAds(t *testing.T) {
	router, mockService := setupAdsAPI(t)

	adsList := []dto.AdResponse{
		{ID: "ad_1", Name: "Test Ad", Status: "ACTIVE"},
	}
	meta := &response.PaginationMeta{Page: 1, Limit: 25, Total: 1, LastPage: 1}

	mockService.On("GetAds", mock.MatchedBy(func(f ads.AdFilter) bool {
		return f.Page == 1 && f.Limit == 25
	})).Return(adsList, meta, nil)

	req, _ := http.NewRequest(http.MethodGet, "/meta/ads?page=1&limit=25", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdsAPI_GetCreative(t *testing.T) {
	router, mockService := setupAdsAPI(t)

	creative := &dto.CreativeResponse{
		ID:   "cr_1",
		Name: "Test Creative",
	}

	mockService.On("GetCreative", "cr_1", ads.DefaultCreativeFields).Return(creative, nil)

	req, _ := http.NewRequest(http.MethodGet, "/meta/creatives/cr_1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
