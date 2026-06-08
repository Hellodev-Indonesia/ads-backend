package integration_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex/ads_backend/internal/core/permission"
	"github.com/alex/ads_backend/internal/core/permission/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupPermissionRouter(mockService permission.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := permission.NewHandler(mockService)
	
	perms := router.Group("/permissions")
	{
		perms.GET("", handler.FindAll)
		perms.GET("/:id", handler.FindByID)
		perms.POST("", handler.Create)
		perms.PUT("/:id", handler.Update)
		perms.DELETE("/:id", handler.Delete)
	}
	return router
}

func TestPermissionAPI_Create(t *testing.T) {
	mockService := permission.NewMockService(t)
	router := setupPermissionRouter(mockService)

	reqPayload := dto.PermissionRequest{
		Name:        "core.user.view",
		Description: "View user",
	}
	body, _ := json.Marshal(reqPayload)

	mockService.On("Create", mock.AnythingOfType("dto.PermissionRequest")).Return(&permission.Permission{
		ID:          1,
		Name:        "core.user.view",
	}, nil)

	req, _ := http.NewRequest(http.MethodPost, "/permissions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code) 
}

func TestPermissionAPI_Update(t *testing.T) {
	mockService := permission.NewMockService(t)
	router := setupPermissionRouter(mockService)

	reqPayload := dto.PermissionRequest{
		Name:        "core.user.edit",
	}
	body, _ := json.Marshal(reqPayload)

	mockService.On("Update", uint(1), reqPayload).Return(&permission.Permission{
		ID:          1,
		Name:        "core.user.edit",
	}, nil)

	req, _ := http.NewRequest(http.MethodPut, "/permissions/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPermissionAPI_Delete(t *testing.T) {
	mockService := permission.NewMockService(t)
	router := setupPermissionRouter(mockService)

	mockService.On("Delete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/permissions/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPermissionAPI_FindByID(t *testing.T) {
	mockService := permission.NewMockService(t)
	router := setupPermissionRouter(mockService)

	mockResp := &dto.PermissionResponse{
		ID:          1,
		Name:        "core.user.view",
	}
	mockService.On("FindByID", uint(1)).Return(mockResp, nil)

	req, _ := http.NewRequest(http.MethodGet, "/permissions/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPermissionAPI_FindByID_NotFound(t *testing.T) {
	mockService := permission.NewMockService(t)
	router := setupPermissionRouter(mockService)

	mockService.On("FindByID", uint(99)).Return((*dto.PermissionResponse)(nil), errors.New("not found"))

	req, _ := http.NewRequest(http.MethodGet, "/permissions/99", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPermissionAPI_FindAll(t *testing.T) {
	mockService := permission.NewMockService(t)
	router := setupPermissionRouter(mockService)

	mockPerms := []dto.PermissionResponse{
		{ID: 1, Name: "core.user.view"},
	}
	mockMeta := &response.PaginationMeta{
		Page:  1,
		Limit: 25,
		Total: 1,
	}

	mockService.On("FindAll", mock.AnythingOfType("dto.PermissionFilter")).Return(mockPerms, mockMeta, nil)

	req, _ := http.NewRequest(http.MethodGet, "/permissions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
