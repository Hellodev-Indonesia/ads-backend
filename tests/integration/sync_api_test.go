package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex/ads_backend/internal/meta/sync"
	"github.com/alex/ads_backend/internal/meta/sync/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupSyncAPI(t *testing.T) (*gin.Engine, *sync.MockJobTrigger) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockJob := sync.NewMockJobTrigger(t)
	// We pass nil for service since we only test endpoints that use JobTrigger
	handler := sync.NewHandler(mockJob, nil)

	syncGroup := router.Group("/meta/sync")
	{
		syncGroup.POST("", handler.TriggerSync)
		syncGroup.GET("/status", handler.SyncStatus)
	}

	return router, mockJob
}

func TestSyncAPI_TriggerSync(t *testing.T) {
	router, mockJob := setupSyncAPI(t)

	reqPayload := dto.TriggerSyncRequest{
		AdAccountID: "act_1",
	}
	body, _ := json.Marshal(reqPayload)

	batches := []*sync.MetaSyncBatch{
		{ID: 1, BatchCode: "B1"},
	}

	mockJob.On("Start", mock.Anything, reqPayload).Return(batches, nil)

	req, _ := http.NewRequest(http.MethodPost, "/meta/sync", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestSyncAPI_TriggerSync_AlreadyRunning(t *testing.T) {
	router, mockJob := setupSyncAPI(t)

	reqPayload := dto.TriggerSyncRequest{}
	body, _ := json.Marshal(reqPayload)

	mockJob.On("Start", mock.Anything, reqPayload).Return(([]*sync.MetaSyncBatch)(nil), sync.ErrAlreadyRunning)

	req, _ := http.NewRequest(http.MethodPost, "/meta/sync", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestSyncAPI_SyncStatus(t *testing.T) {
	router, mockJob := setupSyncAPI(t)

	mockJob.On("IsRunning").Return(true)

	req, _ := http.NewRequest(http.MethodGet, "/meta/sync/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
