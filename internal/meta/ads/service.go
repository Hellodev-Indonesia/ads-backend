package ads

import (
	"encoding/json"
	"net/url"

	"github.com/alex/ads_backend/internal/meta/ads/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
	"github.com/alex/ads_backend/pkg/response"
)

const DefaultAdFields = "id,campaign_id,adset_id,name,status,effective_status,creative,created_time,updated_time"
const DefaultCreativeFields = "id,name,title,body,image_url,thumbnail_url,object_story_spec,asset_feed_spec,url_tags"

type Service interface {
	GetAds(adAccountID string, fields string, limit string, after string, before string, autoPage bool) ([]dto.AdResponse, *response.MetaPaging, error)
	GetCreative(creativeID string, fields string) (*dto.CreativeResponse, error)
}

type serviceImpl struct {
	client *meta_client.Client
}

func NewService(client *meta_client.Client) Service {
	return &serviceImpl{client}
}

func (s *serviceImpl) GetAds(adAccountID string, fields string, limit string, after string, before string, autoPage bool) ([]dto.AdResponse, *response.MetaPaging, error) {
	params := url.Values{}
	params.Set("fields", fields)
	if limit != "" {
		params.Set("limit", limit)
	}
	if after != "" {
		params.Set("after", after)
	}
	if before != "" {
		params.Set("before", before)
	}

	rawList, paging, err := s.client.Get(adAccountID+"/ads", params, autoPage)
	if err != nil {
		return nil, nil, err
	}

	var result []dto.AdResponse
	for _, raw := range rawList {
		var item dto.AdResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			return nil, nil, err
		}
		result = append(result, item)
	}

	return result, mapPaging(paging), nil
}

func (s *serviceImpl) GetCreative(creativeID string, fields string) (*dto.CreativeResponse, error) {
	params := url.Values{}
	params.Set("fields", fields)

	rawList, _, err := s.client.Get(creativeID, params, false)
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

func mapPaging(p *meta_client.Paging) *response.MetaPaging {
	if p == nil {
		return nil
	}
	res := &response.MetaPaging{}
	res.Cursors.Before = p.Cursors.Before
	res.Cursors.After = p.Cursors.After
	res.HasPrevious = p.Previous != ""
	res.HasNext = p.Next != ""
	return res
}
