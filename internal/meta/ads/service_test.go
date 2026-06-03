package ads_test

import (
	"errors"
	"testing"

	"github.com/alex/ads_backend/internal/meta/ads"
	"github.com/stretchr/testify/assert"
)

func setupService(t *testing.T) (*ads.MockRepository, ads.Service) {
	mockRepo := ads.NewMockRepository(t)
	svc := ads.NewService(nil, mockRepo)
	return mockRepo, svc
}

func TestService_GetAds(t *testing.T) {
	mockRepo, svc := setupService(t)

	filter := ads.AdFilter{
		Page:  1,
		Limit: 10,
	}

	adsList := []ads.MetaAd{
		{ID: "ad_1", Name: "Test Ad", Status: "ACTIVE"},
	}

	mockRepo.On("FindAll", filter).Return(adsList, int64(1), nil)

	resp, meta, err := svc.GetAds(filter)

	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "ad_1", resp[0].ID)
	assert.Equal(t, "Test Ad", resp[0].Name)
	assert.NotNil(t, meta)
	assert.Equal(t, 1, meta.Total)
}

func TestService_GetCreative(t *testing.T) {
	mockRepo, svc := setupService(t)

	jsonPayload := `{"id":"cr_1","name":"Test Creative"}`
	mockRepo.On("FindCreativeRawPayload", "cr_1").Return(jsonPayload, nil)

	resp, err := svc.GetCreative("cr_1", ads.DefaultCreativeFields)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "cr_1", resp.ID)
	assert.Equal(t, "Test Creative", resp.Name)
}

func TestService_GetCreative_NotFound(t *testing.T) {
	mockRepo, svc := setupService(t)

	mockRepo.On("FindCreativeRawPayload", "unknown").Return("", errors.New("not found"))

	resp, err := svc.GetCreative("unknown", ads.DefaultCreativeFields)

	assert.Error(t, err)
	assert.Nil(t, resp)
}
