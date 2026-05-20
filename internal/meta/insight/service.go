package insight

import (
	"encoding/json"
	"net/url"

	"github.com/alex/ads_backend/internal/meta/insight/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
)

type Service interface {
	GetInsights(adAccountID string) ([]dto.InsightResponse, error)
}

type serviceImpl struct {
	client *meta_client.Client
}

func NewService(client *meta_client.Client) Service {
	return &serviceImpl{client}
}

func (s *serviceImpl) GetInsights(adAccountID string) ([]dto.InsightResponse, error) {
	params := url.Values{}
	params.Set("fields", "campaign_id,campaign_name,impressions,clicks,spend,cpc,cpm,ctr")
	params.Set("date_preset", "today")

	rawList, err := s.client.Get(adAccountID+"/insights", params, true)
	if err != nil {
		return nil, err
	}

	var result []dto.InsightResponse
	for _, raw := range rawList {
		var item dto.InsightResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}
