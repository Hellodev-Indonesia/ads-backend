package campaign

import (
	"encoding/json"
	"net/url"

	"github.com/alex/ads_backend/internal/meta/campaign/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
	"github.com/alex/ads_backend/pkg/response"
)

const DefaultFields = "id,name,status,effective_status,objective,daily_budget,lifetime_budget,budget_remaining,bid_strategy,start_time,stop_time,created_time,updated_time"

type Service interface {
	GetCampaigns(adAccountID string, fields string, limit string, after string, before string, autoPage bool) ([]dto.CampaignResponse, *response.MetaPaging, error)
}

type serviceImpl struct {
	client *meta_client.Client
}

func NewService(client *meta_client.Client) Service {
	return &serviceImpl{client}
}

func (s *serviceImpl) GetCampaigns(adAccountID string, fields string, limit string, after string, before string, autoPage bool) ([]dto.CampaignResponse, *response.MetaPaging, error) {
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

	rawList, paging, err := s.client.Get(adAccountID+"/campaigns", params, autoPage)
	if err != nil {
		return nil, nil, err
	}

	var result []dto.CampaignResponse
	for _, raw := range rawList {
		var item dto.CampaignResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			return nil, nil, err
		}
		result = append(result, item)
	}

	return result, mapPaging(paging), nil
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
