package campaign

import (
	"encoding/json"
	"net/url"

	"github.com/alex/ads_backend/internal/meta/campaign/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
)

type Service interface {
	GetCampaigns(adAccountID string) ([]dto.CampaignResponse, error)
}

type serviceImpl struct {
	client *meta_client.Client
}

func NewService(client *meta_client.Client) Service {
	return &serviceImpl{client}
}

func (s *serviceImpl) GetCampaigns(adAccountID string) ([]dto.CampaignResponse, error) {
	params := url.Values{}
	params.Set("fields", "id,name,status,effective_status,objective,created_time,updated_time")

	rawList, err := s.client.Get(adAccountID+"/campaigns", params, true)
	if err != nil {
		return nil, err
	}

	var result []dto.CampaignResponse
	for _, raw := range rawList {
		var item dto.CampaignResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}
