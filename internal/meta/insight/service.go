package insight

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/alex/ads_backend/internal/meta/insight/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
	"github.com/alex/ads_backend/pkg/response"
)

const CampaignInsightFields = "account_id,account_name,account_currency,campaign_id,campaign_name,objective,spend,impressions,reach,clicks,inline_link_clicks,inline_link_click_ctr,cpc,cpm,ctr,actions,action_values,cost_per_action_type,date_start,date_stop"
const AdInsightFields = "account_id,account_name,account_currency,campaign_id,campaign_name,adset_id,adset_name,ad_id,ad_name,objective,spend,impressions,reach,clicks,inline_link_clicks,inline_link_click_ctr,cpc,cpm,ctr,actions,action_values,cost_per_action_type,date_start,date_stop"

type Service interface {
	// DB reads (used by handlers)
	GetCampaignInsights(filter InsightFilter) ([]dto.InsightResponse, *response.PaginationMeta, error)
	GetAdInsights(filter InsightFilter) ([]dto.InsightResponse, *response.PaginationMeta, error)

	// Meta API sync (used by sync job)
	SyncCampaignInsights(adAccountID string) (int, error)
	SyncAdInsights(adAccountID string) (int, error)
}

type serviceImpl struct {
	client *meta_client.Client
	repo   Repository
}

func NewService(client *meta_client.Client, repo Repository) Service {
	return &serviceImpl{client: client, repo: repo}
}

// --- DB READ METHODS ---

func (s *serviceImpl) GetCampaignInsights(filter InsightFilter) ([]dto.InsightResponse, *response.PaginationMeta, error) {
	insights, total, err := s.repo.FindCampaignInsights(filter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch campaign insights: %w", err)
	}
	return s.buildInsightResponse(insights, total, filter)
}

func (s *serviceImpl) GetAdInsights(filter InsightFilter) ([]dto.InsightResponse, *response.PaginationMeta, error) {
	insights, total, err := s.repo.FindAdInsights(filter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch ad insights: %w", err)
	}
	return s.buildInsightResponse(insights, total, filter)
}

func (s *serviceImpl) buildInsightResponse(insights []MetaInsight, total int64, filter InsightFilter) ([]dto.InsightResponse, *response.PaginationMeta, error) {
	var result []dto.InsightResponse
	for _, i := range insights {
		result = append(result, mapModelToDTO(i))
	}

	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	lastPage := int(total) / filter.Limit
	if int(total)%filter.Limit > 0 {
		lastPage++
	}

	meta := &response.PaginationMeta{
		Page:     filter.Page,
		Limit:    filter.Limit,
		Total:    int(total),
		LastPage: lastPage,
	}

	return result, meta, nil
}

// --- META API SYNC ---

func (s *serviceImpl) SyncCampaignInsights(adAccountID string) (int, error) {
	return s.syncInsightsInternal(adAccountID, "campaign", CampaignInsightFields)
}

func (s *serviceImpl) SyncAdInsights(adAccountID string) (int, error) {
	return s.syncInsightsInternal(adAccountID, "ad", AdInsightFields)
}

func (s *serviceImpl) syncInsightsInternal(adAccountID string, level string, fields string) (int, error) {
	params := url.Values{}
	params.Set("level", level)
	params.Set("fields", fields)
	params.Set("date_preset", "last_30d")

	rawList, _, err := s.client.Get(adAccountID+"/insights", params, true)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch %s insights from Meta: %w", level, err)
	}

	var models []MetaInsight
	for _, raw := range rawList {
		var item dto.InsightResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			log.Printf("Warning: skipping insight unmarshal error: %v", err)
			continue
		}
		models = append(models, mapDTOToModel(item, level))
	}

	if err := s.repo.UpsertBatch(models); err != nil {
		return 0, fmt.Errorf("failed to upsert %s insights: %w", level, err)
	}

	return len(models), nil
}

// --- MAPPERS ---

func mapModelToDTO(m MetaInsight) dto.InsightResponse {
	resp := dto.InsightResponse{
		AccountID:          m.AccountID,
		AccountName:        m.AccountName,
		AccountCurrency:    m.AccountCurrency,
		CampaignID:         m.CampaignID,
		CampaignName:       m.CampaignName,
		AdSetID:            m.AdSetID,
		AdSetName:          m.AdSetName,
		AdID:               m.AdID,
		AdName:             m.AdName,
		Objective:          m.Objective,
		Impressions:        strconv.FormatInt(m.Impressions, 10),
		Reach:              strconv.FormatInt(m.Reach, 10),
		Clicks:             strconv.FormatInt(m.Clicks, 10),
		InlineLinkClicks:   strconv.FormatInt(m.InlineLinkClicks, 10),
		InlineLinkClickCtr: formatDecimal(m.InlineLinkClickCtr),
		Spend:              formatDecimal(m.Spend),
		CPC:                formatDecimal(m.CPC),
		CPM:                formatDecimal(m.CPM),
		CTR:                formatDecimal(m.CTR),
		DateStart:          m.DateStart,
		DateStop:           m.DateStop,
	}

	// Convert json.RawMessage back to interface{} for DTO
	if m.Actions != nil {
		var v interface{}
		if err := json.Unmarshal(m.Actions, &v); err == nil {
			resp.Actions = v
		}
	}
	if m.ActionValues != nil {
		var v interface{}
		if err := json.Unmarshal(m.ActionValues, &v); err == nil {
			resp.ActionValues = v
		}
	}
	if m.CostPerActionType != nil {
		var v interface{}
		if err := json.Unmarshal(m.CostPerActionType, &v); err == nil {
			resp.CostPerActionType = v
		}
	}

	return resp
}

func mapDTOToModel(d dto.InsightResponse, level string) MetaInsight {
	model := MetaInsight{
		AccountID:          d.AccountID,
		AccountName:        d.AccountName,
		AccountCurrency:    d.AccountCurrency,
		CampaignID:         d.CampaignID,
		CampaignName:       d.CampaignName,
		AdSetID:            d.AdSetID,
		AdSetName:          d.AdSetName,
		AdID:               d.AdID,
		AdName:             d.AdName,
		Level:              level,
		Objective:          d.Objective,
		Impressions:        parseInt(d.Impressions),
		Reach:              parseInt(d.Reach),
		Clicks:             parseInt(d.Clicks),
		InlineLinkClicks:   parseInt(d.InlineLinkClicks),
		InlineLinkClickCtr: parseDecimal(d.InlineLinkClickCtr),
		Spend:              parseDecimal(d.Spend),
		CPC:                parseDecimal(d.CPC),
		CPM:                parseDecimal(d.CPM),
		CTR:                parseDecimal(d.CTR),
		DateStart:          d.DateStart,
		DateStop:           d.DateStop,
	}

	// Convert interface{} to json.RawMessage for storage
	if d.Actions != nil {
		if b, err := json.Marshal(d.Actions); err == nil {
			model.Actions = b
		}
	}
	if d.ActionValues != nil {
		if b, err := json.Marshal(d.ActionValues); err == nil {
			model.ActionValues = b
		}
	}
	if d.CostPerActionType != nil {
		if b, err := json.Marshal(d.CostPerActionType); err == nil {
			model.CostPerActionType = b
		}
	}

	return model
}

func parseInt(s string) int64 {
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func parseDecimal(s string) float64 {
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func formatDecimal(v float64) string {
	if v == 0 {
		return "0"
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}
