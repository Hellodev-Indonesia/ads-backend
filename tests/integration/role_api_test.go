package integration_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex/ads_backend/internal/core/role"
	"github.com/alex/ads_backend/internal/core/role/dto"
	permDto "github.com/alex/ads_backend/internal/core/permission/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(mockService role.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := role.NewHandler(mockService)
	
	roles := router.Group("/roles")
	{
		roles.GET("", handler.FindAll)
		roles.GET("/:id", handler.FindByID)
		roles.POST("", handler.Create)
		roles.PUT("/:id", handler.Update)
		roles.DELETE("/:id", handler.Delete)
		roles.POST("/:id/permissions", handler.AssignPermissions)
	}
	return router
}

func TestRoleAPI_Create(t *testing.T) {
	mockService := role.NewMockService(t)
	router := setupRouter(mockService)

	reqPayload := dto.RoleRequest{
		Name:        "Manager",
		Description: "Manager role",
	}
	body, _ := json.Marshal(reqPayload)

	mockService.On("Create", mock.AnythingOfType("dto.RoleRequest")).Return(&role.Role{
		ID:          2,
		Name:        "Manager",
		Description: "Manager role",
	}, nil)

	req, _ := http.NewRequest(http.MethodPost, "/roles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code) 

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])
	assert.Equal(t, "Role created successfully", resp["message"])
}

func TestRoleAPI_Update(t *testing.T) {
	mockService := role.NewMockService(t)
	router := setupRouter(mockService)

	reqPayload := dto.RoleRequest{
		Name:        "Updated Manager",
		Description: "Updated role",
	}
	body, _ := json.Marshal(reqPayload)

	mockService.On("Update", uint(1), reqPayload).Return(&role.Role{
		ID:          1,
		Name:        "Updated Manager",
	}, nil)

	req, _ := http.NewRequest(http.MethodPut, "/roles/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleAPI_Delete(t *testing.T) {
	mockService := role.NewMockService(t)
	router := setupRouter(mockService)

	mockService.On("Delete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/roles/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleAPI_FindByID(t *testing.T) {
	mockService := role.NewMockService(t)
	router := setupRouter(mockService)

	mockResp := &dto.RoleResponse{
		ID:          1,
		Name:        "Admin",
		Description: "Admin role",
		Permissions: []permDto.PermissionResponse{
			{ID: 1, Name: "core.user.view"},
		},
	}
	mockService.On("FindByID", uint(1)).Return(mockResp, nil)

	req, _ := http.NewRequest(http.MethodGet, "/roles/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleAPI_FindByID_NotFound(t *testing.T) {
	mockService := role.NewMockService(t)
	router := setupRouter(mockService)

	mockService.On("FindByID", uint(99)).Return((*dto.RoleResponse)(nil), errors.New("not found"))

	req, _ := http.NewRequest(http.MethodGet, "/roles/99", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRoleAPI_FindAll(t *testing.T) {
	mockService := role.NewMockService(t)
	router := setupRouter(mockService)

	mockRoles := []dto.RoleResponse{
		{ID: 1, Name: "Admin"},
		{ID: 2, Name: "User"},
	}
	mockMeta := &response.PaginationMeta{
		Page:  1,
		Limit: 25,
		Total: 2,
	}

	mockService.On("FindAll", mock.AnythingOfType("dto.RoleFilter")).Return(mockRoles, mockMeta, nil)

	req, _ := http.NewRequest(http.MethodGet, "/roles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleAPI_AssignPermissions(t *testing.T) {
	mockService := role.NewMockService(t)
	router := setupRouter(mockService)

	reqPayload := dto.AssignPermissionRequest{
		PermissionIDs: []uint{1, 2, 3},
	}
	body, _ := json.Marshal(reqPayload)

	mockService.On("AssignPermissions", uint(1), reqPayload).Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "/roles/1/permissions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
