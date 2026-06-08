package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex/ads_backend/internal/core/brand_whitelist_rule"
	"github.com/alex/ads_backend/internal/core/brand_whitelist_rule/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupWhitelistRouter(mockService brand_whitelist_rule.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := brand_whitelist_rule.NewHandler(mockService)

	rules := router.Group("/brands/:brand_id/whitelist-rules")
	{
		rules.GET("", handler.FindAll)
		rules.GET("/:id", handler.FindByID)
		rules.POST("", handler.Create)
		rules.PUT("/:id", handler.Update)
		rules.DELETE("/:id", handler.Delete)
		rules.POST("/check-url", handler.CheckURL)
	}
	return router
}

func TestWhitelistAPI_Create(t *testing.T) {
	mockService := brand_whitelist_rule.NewMockService(t)
	router := setupWhitelistRouter(mockService)

	isActive := true
	reqPayload := dto.CreateWhitelistRuleRequest{
		Scope:     "domain",
		MatchType: "domain",
		Value:     "example.com",
		IsActive:  &isActive,
	}
	body, _ := json.Marshal(reqPayload)

	mockService.On("Create", uint64(1), mock.AnythingOfType("dto.CreateWhitelistRuleRequest")).Return(dto.WhitelistRuleResponse{
		ID:    1,
		Value: "example.com",
	}, nil)

	req, _ := http.NewRequest(http.MethodPost, "/brands/1/whitelist-rules", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code) 
}

func TestWhitelistAPI_FindAll(t *testing.T) {
	mockService := brand_whitelist_rule.NewMockService(t)
	router := setupWhitelistRouter(mockService)

	mockService.On("FindAll", uint64(1), mock.AnythingOfType("dto.WhitelistRuleFilter")).Return([]dto.WhitelistRuleResponse{{ID: 1, Value: "example.com"}}, int64(1), nil)

	req, _ := http.NewRequest(http.MethodGet, "/brands/1/whitelist-rules", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWhitelistAPI_FindByID(t *testing.T) {
	mockService := brand_whitelist_rule.NewMockService(t)
	router := setupWhitelistRouter(mockService)

	mockService.On("FindByID", uint64(1), uint64(1)).Return(dto.WhitelistRuleResponse{ID: 1, Value: "example.com"}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/brands/1/whitelist-rules/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWhitelistAPI_Update(t *testing.T) {
	mockService := brand_whitelist_rule.NewMockService(t)
	router := setupWhitelistRouter(mockService)

	newValue := "updated.com"
	reqPayload := dto.UpdateWhitelistRuleRequest{
		Value: &newValue,
	}
	body, _ := json.Marshal(reqPayload)

	mockService.On("Update", uint64(1), uint64(1), mock.AnythingOfType("dto.UpdateWhitelistRuleRequest")).Return(dto.WhitelistRuleResponse{
		ID:    1,
		Value: "updated.com",
	}, nil)

	req, _ := http.NewRequest(http.MethodPut, "/brands/1/whitelist-rules/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWhitelistAPI_Delete(t *testing.T) {
	mockService := brand_whitelist_rule.NewMockService(t)
	router := setupWhitelistRouter(mockService)

	mockService.On("Delete", uint64(1), uint64(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/brands/1/whitelist-rules/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWhitelistAPI_CheckURL(t *testing.T) {
	mockService := brand_whitelist_rule.NewMockService(t)
	router := setupWhitelistRouter(mockService)

	reqPayload := dto.CheckURLRequest{
		URL:   "https://example.com",
		Scope: "domain",
	}
	body, _ := json.Marshal(reqPayload)

	mockResult := brand_whitelist_rule.MatchResult{
		Allowed: true,
	}
	mockService.On("IsURLAllowed", uint64(1), "https://example.com", "domain").Return(mockResult, nil)

	req, _ := http.NewRequest(http.MethodPost, "/brands/1/whitelist-rules/check-url", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
