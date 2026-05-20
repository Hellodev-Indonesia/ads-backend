package adset

import (
	"encoding/json"
	"net/url"

	"github.com/alex/ads_backend/internal/meta/adset/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
)

type Service interface {
	GetAdSets(adAccountID string) ([]dto.AdSetResponse, error)
}

type serviceImpl struct {
	client *meta_client.Client
}

func NewService(client *meta_client.Client) Service {
	return &serviceImpl{client}
}

func (s *serviceImpl) GetAdSets(adAccountID string) ([]dto.AdSetResponse, error) {
	params := url.Values{}
	params.Set("fields", "id,name,campaign_id,status,effective_status,daily_budget,lifetime_budget,start_time,end_time")

	rawList, err := s.client.Get(adAccountID+"/adsets", params, true)
	if err != nil {
		return nil, err
	}

	var result []dto.AdSetResponse
	for _, raw := range rawList {
		var item dto.AdSetResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}
