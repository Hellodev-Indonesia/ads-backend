package adset_test

import (
	"testing"

	"github.com/alex/ads_backend/internal/meta/adset"
	"github.com/stretchr/testify/assert"
)

func setupService(t *testing.T) (*adset.MockRepository, adset.Service) {
	mockRepo := adset.NewMockRepository(t)
	svc := adset.NewService(nil, mockRepo)
	return mockRepo, svc
}

func TestService_GetAdSets(t *testing.T) {
	mockRepo, svc := setupService(t)

	filter := adset.AdSetFilter{
		Page:  1,
		Limit: 10,
	}

	adsets := []adset.MetaAdSet{
		{ID: "adset_1", Name: "Test Adset", Status: "ACTIVE"},
	}

	mockRepo.On("FindAll", filter).Return(adsets, int64(1), nil)

	resp, meta, err := svc.GetAdSets(filter)

	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "adset_1", resp[0].ID)
	assert.Equal(t, "Test Adset", resp[0].Name)
	assert.NotNil(t, meta)
	assert.Equal(t, 1, meta.Total)
}
