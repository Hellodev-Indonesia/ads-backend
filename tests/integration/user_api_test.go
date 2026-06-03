package integration_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex/ads_backend/internal/core/user"
	"github.com/alex/ads_backend/internal/core/user/dto"
	"github.com/alex/ads_backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUserRouter(mockService user.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := user.NewHandler(mockService)
	
	users := router.Group("/users")
	{
		users.GET("", handler.FindAll)
		users.GET("/:id", handler.FindByID)
		users.POST("", handler.Create)
		users.PUT("/:id", handler.Update)
		users.DELETE("/:id", handler.Delete)
	}
	return router
}

func TestUserAPI_Create(t *testing.T) {
	mockService := user.NewMockService(t)
	router := setupUserRouter(mockService)

	reqPayload := dto.UserRequest{
		Name:     "John",
		Email:    "john@test.com",
		Password: "password",
	}
	body, _ := json.Marshal(reqPayload)

	mockService.On("Create", mock.AnythingOfType("dto.UserRequest")).Return(&user.User{
		ID:   1,
		Name: "John",
	}, nil)

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code) 
}

func TestUserAPI_Update(t *testing.T) {
	mockService := user.NewMockService(t)
	router := setupUserRouter(mockService)

	reqPayload := dto.UserRequest{
		Name:  "John Updated",
		Email: "john@test.com",
	}
	body, _ := json.Marshal(reqPayload)

	mockService.On("Update", uint(1), reqPayload).Return(&user.User{
		ID:   1,
		Name: "John Updated",
	}, nil)

	req, _ := http.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserAPI_Delete(t *testing.T) {
	mockService := user.NewMockService(t)
	router := setupUserRouter(mockService)

	mockService.On("Delete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserAPI_FindByID(t *testing.T) {
	mockService := user.NewMockService(t)
	router := setupUserRouter(mockService)

	mockResp := &dto.UserResponse{
		ID:    1,
		Name:  "John",
		Email: "john@test.com",
	}
	mockService.On("FindByID", uint(1)).Return(mockResp, nil)

	req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserAPI_FindByID_NotFound(t *testing.T) {
	mockService := user.NewMockService(t)
	router := setupUserRouter(mockService)

	mockService.On("FindByID", uint(99)).Return((*dto.UserResponse)(nil), errors.New("not found"))

	req, _ := http.NewRequest(http.MethodGet, "/users/99", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserAPI_FindAll(t *testing.T) {
	mockService := user.NewMockService(t)
	router := setupUserRouter(mockService)

	mockUsers := []dto.UserResponse{
		{ID: 1, Name: "John"},
	}
	mockMeta := &response.PaginationMeta{
		Page:  1,
		Limit: 25,
		Total: 1,
	}

	mockService.On("FindAll", mock.AnythingOfType("dto.UserFilter")).Return(mockUsers, mockMeta, nil)

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
