package insight_test

import (
	"testing"

	"github.com/alex/ads_backend/internal/meta/insight"
	"github.com/stretchr/testify/assert"
)

func setupService(t *testing.T) (*insight.MockRepository, insight.Service) {
	mockRepo := insight.NewMockRepository(t)
	svc := insight.NewService(nil, mockRepo)
	return mockRepo, svc
}

func TestService_GetCampaignInsights(t *testing.T) {
	mockRepo, svc := setupService(t)

	filter := insight.InsightFilter{
		Page:  1,
		Limit: 10,
	}

	insights := []insight.MetaInsight{
		{CampaignID: "cmp_1", Spend: 100.00},
	}

	mockRepo.On("FindCampaignInsights", filter).Return(insights, int64(1), nil)

	resp, meta, err := svc.GetCampaignInsights(filter)

	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "cmp_1", resp[0].CampaignID)
	assert.Equal(t, "100", resp[0].Spend)
	assert.NotNil(t, meta)
	assert.Equal(t, 1, meta.Total)
}

func TestService_GetAdInsights(t *testing.T) {
	mockRepo, svc := setupService(t)

	filter := insight.InsightFilter{
		Page:  1,
		Limit: 10,
	}

	insights := []insight.MetaInsight{
		{AdID: "ad_1", Spend: 50.00},
	}

	mockRepo.On("FindAdInsights", filter).Return(insights, int64(1), nil)

	resp, meta, err := svc.GetAdInsights(filter)

	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "ad_1", resp[0].AdID)
	assert.Equal(t, "50", resp[0].Spend)
	assert.NotNil(t, meta)
	assert.Equal(t, 1, meta.Total)
}
