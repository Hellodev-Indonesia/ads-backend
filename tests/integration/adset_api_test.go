package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex/ads_backend/internal/meta/adset"
	"github.com/alex/ads_backend/internal/meta/adset/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAdsetAPI(t *testing.T) (*gin.Engine, *adset.MockService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := adset.NewMockService(t)
	handler := adset.NewHandler(mockService)

	adsetGroup := router.Group("/meta/adsets")
	{
		adsetGroup.GET("", handler.GetAdSets)
	}

	return router, mockService
}

func TestAdsetAPI_GetAdsets(t *testing.T) {
	router, mockService := setupAdsetAPI(t)

	adsets := []dto.AdSetResponse{
		{ID: "adset_1", Name: "Test Adset", Status: "ACTIVE"},
	}
	meta := &response.PaginationMeta{Page: 1, Limit: 25, Total: 1, LastPage: 1}

	mockService.On("GetAdSets", mock.MatchedBy(func(f adset.AdSetFilter) bool {
		return f.Page == 1 && f.Limit == 25
	})).Return(adsets, meta, nil)

	req, _ := http.NewRequest(http.MethodGet, "/meta/adsets?page=1&limit=25", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
