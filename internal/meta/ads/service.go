package ads

import (
	"encoding/json"
	"net/url"

	"github.com/alex/ads_backend/internal/meta/ads/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
)

type Service interface {
	GetAds(adAccountID string) ([]dto.AdResponse, error)
	GetCreative(creativeID string) (*dto.CreativeResponse, error)
}

type serviceImpl struct {
	client *meta_client.Client
}

func NewService(client *meta_client.Client) Service {
	return &serviceImpl{client}
}

func (s *serviceImpl) GetAds(adAccountID string) ([]dto.AdResponse, error) {
	params := url.Values{}
	params.Set("fields", "id,name,adset_id,campaign_id,status,effective_status,creative")

	rawList, err := s.client.Get(adAccountID+"/ads", params, true)
	if err != nil {
		return nil, err
	}

	var result []dto.AdResponse
	for _, raw := range rawList {
		var item dto.AdResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (s *serviceImpl) GetCreative(creativeID string) (*dto.CreativeResponse, error) {
	params := url.Values{}
	params.Set("fields", "id,name,title,body,image_url,thumbnail_url,object_story_spec,asset_feed_spec")

	rawList, err := s.client.Get(creativeID, params, false)
	if err != nil {
		return nil, err
	}

	if len(rawList) == 0 {
		return nil, nil
	}

	var item dto.CreativeResponse
	if err := json.Unmarshal(rawList[0], &item); err != nil {
		return nil, err
	}

	return &item, nil
}
