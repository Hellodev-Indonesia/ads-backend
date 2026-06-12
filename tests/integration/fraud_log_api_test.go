package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex/ads_backend/internal/core/fraud_log"
	"github.com/alex/ads_backend/internal/core/fraud_log/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupFraudLogRouter(mockService fraud_log.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := fraud_log.NewHandler(mockService)

	logs := router.Group("/fraud-logs")
	logs.Use(func(c *gin.Context) {
		c.Set("user_id", float64(2)) // float64 because JWT parses numbers as float64
		c.Next()
	})
	{
		logs.GET("", handler.FindAll)
		logs.GET("/:id", handler.FindByID)
		logs.PUT("/:id/resolve", handler.Resolve)
	}
	return router
}

func TestFraudLogAPI_FindAll(t *testing.T) {
	mockService := fraud_log.NewMockService(t)
	router := setupFraudLogRouter(mockService)

	mockService.On("FindAll", mock.AnythingOfType("dto.FraudLogFilter")).Return([]dto.FraudLogResponse{{ID: 1, EventType: "fake_click"}}, int64(1), nil)

	req, _ := http.NewRequest(http.MethodGet, "/fraud-logs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFraudLogAPI_FindByID(t *testing.T) {
	mockService := fraud_log.NewMockService(t)
	router := setupFraudLogRouter(mockService)

	mockService.On("FindByID", uint64(1)).Return(dto.FraudLogResponse{ID: 1, EventType: "fake_click"}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/fraud-logs/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFraudLogAPI_Resolve(t *testing.T) {
	mockService := fraud_log.NewMockService(t)
	router := setupFraudLogRouter(mockService)

	mockService.On("Resolve", uint64(1), uint64(2)).Return(dto.FraudLogResponse{
		ID:     1,
		Status: "resolved",
	}, nil)

	req, _ := http.NewRequest(http.MethodPut, "/fraud-logs/1/resolve", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
