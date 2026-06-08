package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex/ads_backend/internal/core/auth"
	"github.com/alex/ads_backend/internal/core/auth/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAuthRouter(mockService auth.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := auth.NewHandler(mockService)
	
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", handler.Login)
		authRoutes.POST("/logout", handler.Logout)
	}
	return router
}

func TestAuthAPI_Login(t *testing.T) {
	mockService := auth.NewMockService(t)
	router := setupAuthRouter(mockService)

	reqPayload := dto.LoginRequest{
		Email:    "test@test.com",
		Password: "password",
	}
	body, _ := json.Marshal(reqPayload)

	mockResp := &dto.LoginResponse{
		Token: "dummy-token",
		User: dto.AuthUserResponse{
			ID:    1,
			Email: "test@test.com",
		},
	}

	mockService.On("Login", reqPayload).Return(mockResp, nil)

	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code) 
}

func TestAuthAPI_Logout(t *testing.T) {
	mockService := auth.NewMockService(t)
	router := setupAuthRouter(mockService)

	req, _ := http.NewRequest(http.MethodPost, "/auth/logout", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
