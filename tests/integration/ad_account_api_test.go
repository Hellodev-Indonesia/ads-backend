package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex/ads_backend/internal/meta/ad_account"
	"github.com/alex/ads_backend/internal/meta/ad_account/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAdAccountAPI(t *testing.T) (*gin.Engine, *ad_account.MockService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := ad_account.NewMockService(t)
	handler := ad_account.NewHandler(mockService)

	adAccountGroup := router.Group("/meta/ad-accounts")
	{
		adAccountGroup.GET("", handler.GetAdAccounts)
		adAccountGroup.GET("/unassigned", handler.GetUnassignedAdAccounts)
		adAccountGroup.PUT("/brand", handler.AssignBrand)
	}

	return router, mockService
}

func TestAdAccountAPI_GetAdAccounts(t *testing.T) {
	router, mockService := setupAdAccountAPI(t)

	accounts := []dto.AdAccountResponse{
		{ID: "act_1", Name: "Test Account", AccountStatus: 1},
	}
	meta := &response.Meta{Page: 1, Limit: 25, Total: 1, LastPage: 1}

	mockService.On("GetAdAccounts", mock.MatchedBy(func(f ad_account.AdAccountFilter) bool {
		return f.Page == 1 && f.Limit == 25
	})).Return(accounts, meta, nil)

	req, _ := http.NewRequest(http.MethodGet, "/meta/ad-accounts?page=1&limit=25", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdAccountAPI_GetUnassignedAdAccounts(t *testing.T) {
	router, mockService := setupAdAccountAPI(t)

	accounts := []dto.AdAccountResponse{
		{ID: "act_2", Name: "Unassigned Account", AccountStatus: 1},
	}
	meta := &response.Meta{Page: 1, Limit: 25, Total: 1, LastPage: 1}

	mockService.On("GetUnassigned", mock.MatchedBy(func(f ad_account.AdAccountFilter) bool {
		return f.Page == 1 && f.Limit == 25
	})).Return(accounts, meta, nil)

	req, _ := http.NewRequest(http.MethodGet, "/meta/ad-accounts/unassigned?page=1&limit=25", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdAccountAPI_AssignBrand(t *testing.T) {
	router, mockService := setupAdAccountAPI(t)

	brandID := uint64(1)
	reqPayload := dto.AssignBrandRequest{
		AdAccountIDs: []string{"act_1", "act_2"},
		BrandID:      &brandID,
	}
	body, _ := json.Marshal(reqPayload)

	mockService.On("BulkAssignBrand", reqPayload).Return(nil)

	req, _ := http.NewRequest(http.MethodPut, "/meta/ad-accounts/brand", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdAccountAPI_UnassignBrand(t *testing.T) {
	router, mockService := setupAdAccountAPI(t)

	reqPayload := dto.AssignBrandRequest{
		AdAccountIDs: []string{"act_1", "act_2"},
		BrandID:      nil,
	}
	body, _ := json.Marshal(reqPayload)

	mockService.On("BulkAssignBrand", reqPayload).Return(nil)

	req, _ := http.NewRequest(http.MethodPut, "/meta/ad-accounts/brand", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
