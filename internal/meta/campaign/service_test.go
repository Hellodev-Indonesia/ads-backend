package campaign_test

import (
	"errors"
	"testing"

	"github.com/alex/ads_backend/internal/meta/campaign"
	"github.com/stretchr/testify/assert"
)

func setupService(t *testing.T) (*campaign.MockRepository, campaign.Service) {
	mockRepo := campaign.NewMockRepository(t)
	svc := campaign.NewService(nil, mockRepo)
	return mockRepo, svc
}

func TestService_GetCampaigns(t *testing.T) {
	mockRepo, svc := setupService(t)

	filter := campaign.CampaignFilter{
		Page:  1,
		Limit: 10,
	}

	campaigns := []campaign.MetaCampaign{
		{ID: "cmp_1", Name: "Test Campaign", Status: "ACTIVE"},
	}

	mockRepo.On("FindAll", filter).Return(campaigns, int64(1), nil)

	resp, meta, err := svc.GetCampaigns(filter)

	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "cmp_1", resp[0].ID)
	assert.Equal(t, "Test Campaign", resp[0].Name)
	assert.NotNil(t, meta)
	assert.Equal(t, 1, meta.Total)
}

func TestService_GetCampaignByID(t *testing.T) {
	mockRepo, svc := setupService(t)

	existing := &campaign.MetaCampaign{
		ID:     "cmp_1",
		Name:   "Test Campaign",
		Status: "ACTIVE",
	}

	mockRepo.On("FindByID", "cmp_1").Return(existing, nil)

	resp, err := svc.GetCampaignByID("cmp_1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "cmp_1", resp.ID)
}

func TestService_GetCampaignByID_NotFound(t *testing.T) {
	mockRepo, svc := setupService(t)

	mockRepo.On("FindByID", "unknown").Return(nil, errors.New("not found"))

	resp, err := svc.GetCampaignByID("unknown")

	assert.Error(t, err)
	assert.Nil(t, resp)
}
