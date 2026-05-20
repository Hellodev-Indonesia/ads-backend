package insight

import (
	"encoding/json"
	"net/url"

	"github.com/alex/ads_backend/internal/meta/insight/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
	"github.com/alex/ads_backend/pkg/response"
)

const CampaignInsightFields = "account_id,account_name,account_currency,campaign_id,campaign_name,objective,spend,impressions,reach,clicks,inline_link_clicks,inline_link_click_ctr,cpc,cpm,ctr,actions,action_values,cost_per_action_type,date_start,date_stop"
const AdInsightFields = "account_id,account_name,account_currency,campaign_id,campaign_name,adset_id,adset_name,ad_id,ad_name,objective,spend,impressions,reach,clicks,inline_link_clicks,inline_link_click_ctr,cpc,cpm,ctr,actions,action_values,cost_per_action_type,date_start,date_stop"

type Service interface {
	GetCampaignInsights(adAccountID string, fields string, limit string, after string, before string, autoPage bool) ([]dto.InsightResponse, *response.MetaPaging, error)
	GetAdInsights(adAccountID string, fields string, limit string, after string, before string, autoPage bool) ([]dto.InsightResponse, *response.MetaPaging, error)
}

type serviceImpl struct {
	client *meta_client.Client
}

func NewService(client *meta_client.Client) Service {
	return &serviceImpl{client}
}

func (s *serviceImpl) getInsightsInternal(adAccountID string, level string, fields string, limit string, after string, before string, autoPage bool) ([]dto.InsightResponse, *response.MetaPaging, error) {
	params := url.Values{}
	params.Set("level", level)
	params.Set("fields", fields)
	params.Set("date_preset", "last_30d") // Default based on user request
	if limit != "" {
		params.Set("limit", limit)
	}
	if after != "" {
		params.Set("after", after)
	}
	if before != "" {
		params.Set("before", before)
	}

	rawList, paging, err := s.client.Get(adAccountID+"/insights", params, autoPage)
	if err != nil {
		return nil, nil, err
	}

	var result []dto.InsightResponse
	for _, raw := range rawList {
		var item dto.InsightResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			return nil, nil, err
		}
		result = append(result, item)
	}

	return result, mapPaging(paging), nil
}

func (s *serviceImpl) GetCampaignInsights(adAccountID string, fields string, limit string, after string, before string, autoPage bool) ([]dto.InsightResponse, *response.MetaPaging, error) {
	return s.getInsightsInternal(adAccountID, "campaign", fields, limit, after, before, autoPage)
}

func (s *serviceImpl) GetAdInsights(adAccountID string, fields string, limit string, after string, before string, autoPage bool) ([]dto.InsightResponse, *response.MetaPaging, error) {
	return s.getInsightsInternal(adAccountID, "ad", fields, limit, after, before, autoPage)
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
