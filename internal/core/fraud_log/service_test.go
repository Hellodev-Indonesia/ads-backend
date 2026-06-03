package fraud_log_test

import (
	"errors"
	"testing"
	"time"

	"github.com/alex/ads_backend/internal/core/fraud_log"
	"github.com/alex/ads_backend/internal/core/fraud_log/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupService(t *testing.T) (*fraud_log.MockRepository, fraud_log.Service) {
	mockRepo := fraud_log.NewMockRepository(t)
	svc := fraud_log.NewService(mockRepo)
	return mockRepo, svc
}

func ptr(i uint64) *uint64 { return &i }

func TestService_Create(t *testing.T) {
	mockRepo, svc := setupService(t)

	req := dto.CreateFraudLogInput{
		BrandID:   ptr(1),
		EventType: "fake_click",
		Severity:  "high",
		Message:   "detected bot",
	}

	mockRepo.On("Create", mock.AnythingOfType("*fraud_log.FraudLog")).Run(func(args mock.Arguments) {
		f := args.Get(0).(*fraud_log.FraudLog)
		f.ID = 1
		f.CreatedAt = time.Now()
		f.UpdatedAt = time.Now()
	}).Return(nil)

	resp, err := svc.Create(req)

	assert.NoError(t, err)
	assert.Equal(t, uint64(1), resp.ID)
	assert.Equal(t, "fake_click", resp.EventType)
	assert.Equal(t, "high", resp.Severity)
	assert.NotNil(t, resp.Message)
	assert.Equal(t, "detected bot", *resp.Message)
}

func TestService_FindByID(t *testing.T) {
	mockRepo, svc := setupService(t)

	now := time.Now()
	msg := "test message"
	existing := &fraud_log.FraudLog{
		ID:         1,
		BrandID:    ptr(1),
		EventType:  "fake_click",
		Message:    &msg,
		CreatedAt:  now,
		UpdatedAt:  now,
		DetectedAt: &now,
	}

	mockRepo.On("FindByID", uint64(1)).Return(existing, nil)

	resp, err := svc.FindByID(1)

	assert.NoError(t, err)
	assert.Equal(t, uint64(1), resp.ID)
	assert.Equal(t, "fake_click", resp.EventType)
}

func TestService_FindAll(t *testing.T) {
	mockRepo, svc := setupService(t)

	filter := dto.FraudLogFilter{
		Page:  1,
		Limit: 10,
	}

	logs := []fraud_log.FraudLog{
		{ID: 1, BrandID: ptr(1), EventType: "fake_click", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	mockRepo.On("FindAll", filter).Return(logs, int64(1), nil)

	resp, total, err := svc.FindAll(filter)

	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, int64(1), total)
}

func TestService_Resolve(t *testing.T) {
	mockRepo, svc := setupService(t)

	now := time.Now()
	existing := &fraud_log.FraudLog{
		ID:        1,
		BrandID:   ptr(1),
		Status:    "open",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("FindByID", uint64(1)).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*fraud_log.FraudLog")).Return(nil)

	resp, err := svc.Resolve(1)

	assert.NoError(t, err)
	assert.Equal(t, "resolved", resp.Status)
}

func TestService_Resolve_AlreadyResolved(t *testing.T) {
	mockRepo, svc := setupService(t)

	now := time.Now()
	existing := &fraud_log.FraudLog{
		ID:        1,
		BrandID:   ptr(1),
		Status:    "resolved",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("FindByID", uint64(1)).Return(existing, nil)

	resp, err := svc.Resolve(1)

	assert.NoError(t, err)
	assert.Equal(t, "resolved", resp.Status)
	// Update shouldn't be called
}

func TestService_ExistsOpenDuplicate(t *testing.T) {
	mockRepo, svc := setupService(t)

	mockRepo.On("ExistsOpenDuplicate", "creative_1", "fake_click", "value").Return(true, nil)

	exists, err := svc.ExistsOpenDuplicate("creative_1", "fake_click", "value")

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestService_FindByID_NotFound(t *testing.T) {
	mockRepo, svc := setupService(t)

	mockRepo.On("FindByID", uint64(99)).Return(nil, errors.New("not found"))

	resp, err := svc.FindByID(99)

	assert.Error(t, err)
	assert.Equal(t, uint64(0), resp.ID)
}
