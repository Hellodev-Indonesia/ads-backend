package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex/ads_backend/internal/core/brand"
	"github.com/alex/ads_backend/internal/core/brand/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupBrandRouter(mockService brand.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := brand.NewHandler(mockService)

	brands := router.Group("/brands")
	{
		brands.GET("", handler.FindAll)
		brands.GET("/:id", handler.FindByID)
		brands.POST("", handler.Create)
		brands.PUT("/:id", handler.Update)
		brands.DELETE("/:id", handler.Delete)
	}
	return router
}

func TestBrandAPI_Create(t *testing.T) {
	mockService := brand.NewMockService(t)
	router := setupBrandRouter(mockService)

	isActive := true
	reqPayload := dto.CreateBrandRequest{
		Name:     "Brand 1",
		IsActive: &isActive,
	}
	body, _ := json.Marshal(reqPayload)

	mockService.On("Create", mock.AnythingOfType("dto.CreateBrandRequest")).Return(dto.BrandResponse{
		ID:   1,
		Name: "Brand 1",
	}, nil)

	req, _ := http.NewRequest(http.MethodPost, "/brands", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code) // success response uses 200
}

func TestBrandAPI_FindAll(t *testing.T) {
	mockService := brand.NewMockService(t)
	router := setupBrandRouter(mockService)

	mockService.On("FindAll", mock.AnythingOfType("dto.BrandFilter")).Return([]dto.BrandResponse{{ID: 1, Name: "Brand 1"}}, int64(1), nil)

	req, _ := http.NewRequest(http.MethodGet, "/brands", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBrandAPI_FindByID(t *testing.T) {
	mockService := brand.NewMockService(t)
	router := setupBrandRouter(mockService)

	mockService.On("FindByID", uint64(1)).Return(dto.BrandResponse{ID: 1, Name: "Brand 1"}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/brands/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBrandAPI_Update(t *testing.T) {
	mockService := brand.NewMockService(t)
	router := setupBrandRouter(mockService)

	newName := "Brand Updated"
	reqPayload := dto.UpdateBrandRequest{
		Name: &newName,
	}
	body, _ := json.Marshal(reqPayload)

	mockService.On("Update", uint64(1), mock.AnythingOfType("dto.UpdateBrandRequest")).Return(dto.BrandResponse{
		ID:   1,
		Name: "Brand Updated",
	}, nil)

	req, _ := http.NewRequest(http.MethodPut, "/brands/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBrandAPI_Delete(t *testing.T) {
	mockService := brand.NewMockService(t)
	router := setupBrandRouter(mockService)

	mockService.On("Delete", uint64(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/brands/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
