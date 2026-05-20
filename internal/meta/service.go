package meta

import (
	"encoding/json"
	"net/url"

	"github.com/alex/ads_backend/internal/meta/dto"
)

type Service interface {
	GetAdAccounts() ([]dto.AdAccountResponse, error)
	GetCampaigns(adAccountID string) ([]dto.CampaignResponse, error)
	GetAdSets(adAccountID string) ([]dto.AdSetResponse, error)
	GetAds(adAccountID string) ([]dto.AdResponse, error)
	GetCreative(creativeID string) (*dto.CreativeResponse, error)
	GetInsights(adAccountID string) ([]dto.InsightResponse, error)
}

type serviceImpl struct {
	client Client
}

func NewService(client Client) Service {
	return &serviceImpl{client}
}

func (s *serviceImpl) GetAdAccounts() ([]dto.AdAccountResponse, error) {
	params := url.Values{}
	params.Set("fields", "id,name,account_status")

	rawList, err := s.client.GetClient().Get("me/adaccounts", params, true)
	if err != nil {
		return nil, err
	}

	var result []dto.AdAccountResponse
	for _, raw := range rawList {
		var item dto.AdAccountResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (s *serviceImpl) GetCampaigns(adAccountID string) ([]dto.CampaignResponse, error) {
	params := url.Values{}
	params.Set("fields", "id,name,status,effective_status,objective,created_time,updated_time")

	rawList, err := s.client.GetClient().Get(adAccountID+"/campaigns", params, true)
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

func (s *serviceImpl) GetAdSets(adAccountID string) ([]dto.AdSetResponse, error) {
	params := url.Values{}
	params.Set("fields", "id,name,campaign_id,status,effective_status,daily_budget,lifetime_budget,start_time,end_time")

	rawList, err := s.client.GetClient().Get(adAccountID+"/adsets", params, true)
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

func (s *serviceImpl) GetAds(adAccountID string) ([]dto.AdResponse, error) {
	params := url.Values{}
	params.Set("fields", "id,name,adset_id,campaign_id,status,effective_status,creative")

	rawList, err := s.client.GetClient().Get(adAccountID+"/ads", params, true)
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

	rawList, err := s.client.GetClient().Get(creativeID, params, false)
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

func (s *serviceImpl) GetInsights(adAccountID string) ([]dto.InsightResponse, error) {
	params := url.Values{}
	params.Set("fields", "campaign_id,campaign_name,impressions,clicks,spend,cpc,cpm,ctr")
	params.Set("date_preset", "today")

	rawList, err := s.client.GetClient().Get(adAccountID+"/insights", params, true)
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
