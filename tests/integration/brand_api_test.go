package integration_test

import (
	"bytes"
	"mime/multipart"
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
		brands.GET("/:slug", handler.FindBySlug)
		brands.POST("", handler.Create)
		brands.PUT("/:slug", handler.Update)
		brands.DELETE("/:slug", handler.Delete)
	}
	return router
}

func TestBrandAPI_Create(t *testing.T) {
	mockService := brand.NewMockService(t)
	router := setupBrandRouter(mockService)

	// removed reqPayload

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "Brand 1")
	writer.WriteField("is_active", "true")
	writer.Close()

	mockService.On("Create", mock.AnythingOfType("dto.CreateBrandRequest")).Return(dto.BrandResponse{
		ID:   1,
		Slug: "brand-1",
		Name: "Brand 1",
	}, nil)

	req, _ := http.NewRequest(http.MethodPost, "/brands", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
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

func TestBrandAPI_FindBySlug(t *testing.T) {
	mockService := brand.NewMockService(t)
	router := setupBrandRouter(mockService)

	mockService.On("FindBySlug", "brand-1").Return(dto.BrandResponse{ID: 1, Slug: "brand-1", Name: "Brand 1"}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/brands/brand-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBrandAPI_Update(t *testing.T) {
	mockService := brand.NewMockService(t)
	router := setupBrandRouter(mockService)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "Brand Updated")
	writer.Close()

	mockService.On("Update", "brand-1", mock.AnythingOfType("dto.UpdateBrandRequest")).Return(dto.BrandResponse{
		ID:   1,
		Slug: "brand-updated",
		Name: "Brand Updated",
	}, nil)

	req, _ := http.NewRequest(http.MethodPut, "/brands/brand-1", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBrandAPI_Delete(t *testing.T) {
	mockService := brand.NewMockService(t)
	router := setupBrandRouter(mockService)

	mockService.On("Delete", "brand-1").Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/brands/brand-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
