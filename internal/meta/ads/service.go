package ads

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/alex/ads_backend/internal/meta/ads/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
	"github.com/alex/ads_backend/pkg/response"
)

// Action type keys used by Meta Graph API
const (
	actionTotalMessaging = "onsite_conversion.total_messaging_connection"
	actionNewMessaging   = "onsite_conversion.messaging_first_reply"
	actionPurchase       = "purchase"
)

const DefaultAdFields = "id,campaign_id,adset_id,name,status,effective_status,creative,created_time,updated_time"
const DefaultCreativeFields = "id,name,title,body,image_url,thumbnail_url,object_story_spec,asset_feed_spec,url_tags"

type Service interface {
	// DB reads (used by handlers)
	GetAds(filter AdFilter) ([]dto.AdResponse, *response.PaginationMeta, error)
	GetAdDashboard(filter AdFilter) ([]dto.AdDashboardRow, *response.PaginationMeta, error)

	// Direct Meta API call (creatives stay as direct calls)
	GetCreative(creativeID string, fields string) (*dto.CreativeResponse, error)

	// Meta API sync (used by sync job)
	SyncAds(adAccountID string) (int, error)
	SyncAdsWithList(adAccountID string) (int, []MetaAd, error)
	SyncAdsByIDs(adAccountID string, ids []string) (int, []MetaAd, error)
}

type serviceImpl struct {
	client *meta_client.Client
	repo   Repository
}

func NewService(client *meta_client.Client, repo Repository) Service {
	return &serviceImpl{client: client, repo: repo}
}

// --- DB READ METHODS ---

func (s *serviceImpl) GetAds(filter AdFilter) ([]dto.AdResponse, *response.PaginationMeta, error) {
	adsList, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch ads: %w", err)
	}

	var result []dto.AdResponse
	for _, a := range adsList {
		result = append(result, mapModelToDTO(a))
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

// --- LOCAL DB READ (creatives are synced in background) ---

func (s *serviceImpl) GetCreative(creativeID string, fields string) (*dto.CreativeResponse, error) {
	// The fields parameter is ignored since we return the exact JSON payload
	// that was synced to the local database in the background job.
	payload, err := s.repo.FindCreativeRawPayload(creativeID)
	if err != nil {
		return nil, fmt.Errorf("creative not found in database: %w", err)
	}

	var item dto.CreativeResponse
	if err := json.Unmarshal([]byte(payload), &item); err != nil {
		return nil, fmt.Errorf("failed to parse creative JSON from db: %w", err)
	}

	return &item, nil
}

// --- META API SYNC ---

func (s *serviceImpl) SyncAds(adAccountID string) (int, error) {
	count, _, err := s.SyncAdsWithList(adAccountID)
	return count, err
}

func (s *serviceImpl) SyncAdsWithList(adAccountID string) (int, []MetaAd, error) {
	params := url.Values{}
	params.Set("fields", DefaultAdFields)

	rawList, _, err := s.client.Get(adAccountID+"/ads", params, true)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to fetch ads from Meta: %w", err)
	}

	var models []MetaAd
	for _, raw := range rawList {
		var item dto.AdResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			log.Printf("Warning: skipping ad unmarshal error: %v", err)
			continue
		}
		models = append(models, mapDTOToModel(item))
	}

	if err := s.repo.UpsertBatch(models); err != nil {
		return 0, nil, fmt.Errorf("failed to upsert ads: %w", err)
	}

	return len(models), models, nil
}

func (s *serviceImpl) SyncAdsByIDs(adAccountID string, ids []string) (int, []MetaAd, error) {
	if len(ids) == 0 {
		return 0, nil, nil
	}

	batchSize := 50
	var allModels []MetaAd

	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}
		batchIDs := ids[i:end]

		params := url.Values{}
		params.Set("ids", strings.Join(batchIDs, ","))
		params.Set("fields", DefaultAdFields)

		rawList, _, err := s.client.Get("", params, false)
		if err != nil {
			log.Printf("Warning: failed to bulk fetch ads: %v", err)
			continue
		}

		if len(rawList) == 0 {
			continue
		}

		var rawMap map[string]json.RawMessage
		if err := json.Unmarshal(rawList[0], &rawMap); err != nil {
			log.Printf("Warning: failed to unmarshal bulk ads response: %v", err)
			continue
		}

		for _, rawPayloadBytes := range rawMap {
			var item dto.AdResponse
			if err := json.Unmarshal(rawPayloadBytes, &item); err != nil {
				continue
			}
			allModels = append(allModels, mapDTOToModel(item))
		}
	}

	if len(allModels) > 0 {
		if err := s.repo.UpsertBatch(allModels); err != nil {
			return 0, nil, fmt.Errorf("failed to upsert ads: %w", err)
		}
	}

	return len(allModels), allModels, nil
}

// --- MAPPERS ---

func mapModelToDTO(m MetaAd) dto.AdResponse {
	return dto.AdResponse{
		ID:              m.ID,
		Name:            m.Name,
		AdSetID:         m.AdSetID,
		CampaignID:      m.CampaignID,
		Status:          m.Status,
		EffectiveStatus: m.EffectiveStatus,
		Creative:        dto.CreativeRef{ID: m.CreativeID},
		CreatedTime:     formatTime(m.CreatedTime),
		UpdatedTime:     formatTime(m.UpdatedTime),
	}
}

func mapDTOToModel(d dto.AdResponse) MetaAd {
	return MetaAd{
		ID:              d.ID,
		CampaignID:      d.CampaignID,
		AdSetID:         d.AdSetID,
		Name:            d.Name,
		Status:          d.Status,
		EffectiveStatus: d.EffectiveStatus,
		CreativeID:      d.Creative.ID,
		CreatedTime:     parseTime(d.CreatedTime),
		UpdatedTime:     parseTime(d.UpdatedTime),
	}
}

func parseTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05-0700",
		"2006-01-02T15:04:05+0700",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return &t
		}
	}
	return nil
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func (s *serviceImpl) GetAdDashboard(filter AdFilter) ([]dto.AdDashboardRow, *response.PaginationMeta, error) {
	rows, total, err := s.repo.FindAdDashboard(filter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch ad dashboard: %w", err)
	}

	result := make([]dto.AdDashboardRow, 0, len(rows))
	for _, r := range rows {
		dtoRow := mapAdScanToDTO(r)
		if dtoRow.DateStart == "" && filter.DateStart != "" {
			dtoRow.DateStart = filter.DateStart
		}
		if dtoRow.DateStop == "" && filter.DateStop != "" {
			dtoRow.DateStop = filter.DateStop
		}
		result = append(result, dtoRow)
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
func mapAdScanToDTO(r adDashboardScan) dto.AdDashboardRow {
	row := dto.AdDashboardRow{
		AdID:            r.AdID,
		AdSetID:         r.AdSetID,
		CampaignID:      r.CampaignID,
		CampaignName:    r.CampaignName,
		AdSetName:       r.AdSetName,
		AdName:          r.AdName,
		Status:          r.Status,
		EffectiveStatus: r.EffectiveStatus,
		CreativeID:      r.CreativeID,
		AmountSpent:     formatNullFloat(r.Spend),
		Impressions:     formatNullInt(r.Impressions),
		Reach:           formatNullInt(r.Reach),
	}

	if r.UpdatedTime != nil {
		row.LastSignificantEdit = r.UpdatedTime.Format(time.RFC3339)
	}

	if r.DateStart != nil && len(*r.DateStart) >= 10 {
		row.DateStart = (*r.DateStart)[:10]
	}
	if r.DateStop != nil && len(*r.DateStop) >= 10 {
		row.DateStop = (*r.DateStop)[:10]
	}

	actions := parseActions(r.Actions)
	row.TotalMessagingConversations = findAction(actions, actionTotalMessaging)
	row.NewMessagingConnections = findAction(actions, actionNewMessaging)
	row.Purchases = findAction(actions, actionPurchase)

	var primaryActionType string
	if len(actions) > 0 {
		row.Results = actions[0].Value
		primaryActionType = actions[0].ActionType
	}

	if val := findAction(actions, actionNewMessaging); val != "0" {
		row.Results = val
		primaryActionType = actionNewMessaging
	} else if val := findAction(actions, actionTotalMessaging); val != "0" {
		row.Results = val
		primaryActionType = actionTotalMessaging
	} else if val := findAction(actions, actionPurchase); val != "0" {
		row.Results = val
		primaryActionType = actionPurchase
	}

	costs := parseActions(r.CostPerActionType)
	if primaryActionType != "" {
		row.CostPerResult = findAction(costs, primaryActionType)
	} else if len(costs) > 0 {
		row.CostPerResult = costs[0].Value
	}

	if row.CostPerResult == "" || row.CostPerResult == "0" {
		spent, _ := strconv.ParseFloat(row.AmountSpent, 64)
		results, _ := strconv.ParseFloat(row.Results, 64)
		if results > 0 {
			row.CostPerResult = formatFloat(math.Ceil(spent / results))
		} else {
			row.CostPerResult = "0"
		}
	} else {
		val, _ := strconv.ParseFloat(row.CostPerResult, 64)
		row.CostPerResult = formatFloat(math.Ceil(val))
	}

	spent, _ := strconv.ParseFloat(row.AmountSpent, 64)
	purchases, _ := strconv.ParseFloat(row.Purchases, 64)
	if purchases > 0 {
		row.CostPerPurchase = formatFloat(math.Ceil(spent / purchases))
	} else {
		row.CostPerPurchase = "0"
	}

	return row
}

type metaAction struct {
	ActionType string `json:"action_type"`
	Value      string `json:"value"`
}

func parseActions(raw json.RawMessage) []metaAction {
	if raw == nil {
		return nil
	}
	var actions []metaAction
	_ = json.Unmarshal(raw, &actions)
	return actions
}

func findAction(actions []metaAction, actionType string) string {
	for _, a := range actions {
		if a.ActionType == actionType {
			return a.Value
		}
	}
	return "0"
}

func formatNullFloat(v *float64) string {
	if v == nil {
		return "0"
	}
	return formatFloat(*v)
}

func formatFloat(v float64) string {
	if v == math.Trunc(v) {
		return fmt.Sprintf("%.0f", v)
	}
	return fmt.Sprintf("%.2f", v)
}

func formatNullInt(v *int64) string {
	if v == nil {
		return "0"
	}
	return strconv.FormatInt(*v, 10)
}
