package adset

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/alex/ads_backend/internal/meta/adset/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
	"github.com/alex/ads_backend/pkg/response"
)

// Action type keys used by Meta Graph API
const (
	actionTotalMessaging = "onsite_conversion.total_messaging_connection"
	actionNewMessaging   = "onsite_conversion.messaging_first_reply"
	actionPurchase       = "purchase"
)

const DefaultFields = "id,name,status,effective_status,campaign_id,daily_budget,lifetime_budget,budget_remaining,bid_strategy,start_time,end_time,created_time,updated_time,targeting"

type Service interface {
	// DB reads (used by handlers)
	GetAdSets(filter AdSetFilter) ([]dto.AdSetResponse, *response.PaginationMeta, error)
	GetAdSetDashboard(filter AdSetFilter) ([]dto.AdSetDashboardRow, *response.PaginationMeta, error)

	// Meta API sync (used by sync job)
	SyncAdSets(adAccountID string) (int, error)
	SyncAdSetsByIDs(adAccountID string, ids []string) (int, error)
}

type serviceImpl struct {
	client *meta_client.Client
	repo   Repository
}

func NewService(client *meta_client.Client, repo Repository) Service {
	return &serviceImpl{client: client, repo: repo}
}

// --- DB READ METHODS ---

func (s *serviceImpl) GetAdSets(filter AdSetFilter) ([]dto.AdSetResponse, *response.PaginationMeta, error) {
	adsets, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch adsets: %w", err)
	}

	var result []dto.AdSetResponse
	for _, a := range adsets {
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

// --- META API SYNC ---

func (s *serviceImpl) SyncAdSets(adAccountID string) (int, error) {
	params := url.Values{}
	params.Set("fields", DefaultFields)

	rawList, _, err := s.client.Get(adAccountID+"/adsets", params, true)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch adsets from Meta: %w", err)
	}

	var models []MetaAdSet
	for _, raw := range rawList {
		var item dto.AdSetResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			log.Printf("Warning: skipping adset unmarshal error: %v", err)
			continue
		}
		models = append(models, mapDTOToModel(item))
	}

	if err := s.repo.UpsertBatch(models); err != nil {
		return 0, fmt.Errorf("failed to upsert adsets: %w", err)
	}

	return len(models), nil
}

func (s *serviceImpl) SyncAdSetsByIDs(adAccountID string, ids []string) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	batchSize := 50
	var allModels []MetaAdSet

	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}
		batchIDs := ids[i:end]

		params := url.Values{}
		params.Set("ids", strings.Join(batchIDs, ","))
		params.Set("fields", DefaultFields)

		rawList, _, err := s.client.Get("", params, false)
		if err != nil {
			log.Printf("Warning: failed to bulk fetch adsets: %v", err)
			continue
		}

		if len(rawList) == 0 {
			continue
		}

		var rawMap map[string]json.RawMessage
		if err := json.Unmarshal(rawList[0], &rawMap); err != nil {
			log.Printf("Warning: failed to unmarshal bulk adsets response: %v", err)
			continue
		}

		for _, rawPayloadBytes := range rawMap {
			var item dto.AdSetResponse
			if err := json.Unmarshal(rawPayloadBytes, &item); err != nil {
				continue
			}
			allModels = append(allModels, mapDTOToModel(item))
		}
	}

	if len(allModels) > 0 {
		if err := s.repo.UpsertBatch(allModels); err != nil {
			return 0, fmt.Errorf("failed to upsert adsets: %w", err)
		}
	}

	return len(allModels), nil
}

// --- MAPPERS ---

func mapModelToDTO(m MetaAdSet) dto.AdSetResponse {
	return dto.AdSetResponse{
		ID:              m.ID,
		Name:            m.Name,
		CampaignID:      m.CampaignID,
		Status:          m.Status,
		EffectiveStatus: m.EffectiveStatus,
		DailyBudget:     formatDecimal(m.DailyBudget),
		LifetimeBudget:  formatDecimal(m.LifetimeBudget),
		BudgetRemaining: formatDecimal(m.BudgetRemaining),
		BidStrategy:     m.BidStrategy,
		AttributionSpec: m.AttributionSpec,
		StartTime:       formatTime(m.StartTime),
		EndTime:         formatTime(m.EndTime),
		CreatedTime:     formatTime(m.CreatedTime),
		UpdatedTime:     formatTime(m.UpdatedTime),
	}
}

func mapDTOToModel(d dto.AdSetResponse) MetaAdSet {
	var attrSpec json.RawMessage
	if d.AttributionSpec != nil {
		attrSpec, _ = json.Marshal(d.AttributionSpec)
	}

	return MetaAdSet{
		ID:              d.ID,
		CampaignID:      d.CampaignID,
		Name:            d.Name,
		Status:          d.Status,
		EffectiveStatus: d.EffectiveStatus,
		DailyBudget:     parseDecimal(d.DailyBudget),
		LifetimeBudget:  parseDecimal(d.LifetimeBudget),
		BudgetRemaining: parseDecimal(d.BudgetRemaining),
		BidStrategy:     d.BidStrategy,
		AttributionSpec: attrSpec,
		StartTime:       parseTime(d.StartTime),
		EndTime:         parseTime(d.EndTime),
		CreatedTime:     parseTime(d.CreatedTime),
		UpdatedTime:     parseTime(d.UpdatedTime),
	}
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

func (s *serviceImpl) GetAdSetDashboard(filter AdSetFilter) ([]dto.AdSetDashboardRow, *response.PaginationMeta, error) {
	rows, total, err := s.repo.FindAdSetDashboard(filter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch adset dashboard: %w", err)
	}

	result := make([]dto.AdSetDashboardRow, 0, len(rows))
	for _, r := range rows {
		dtoRow := mapAdSetScanToDTO(r)
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
func mapAdSetScanToDTO(r adSetDashboardScan) dto.AdSetDashboardRow {
	row := dto.AdSetDashboardRow{
		AdSetID:         r.AdSetID,
		CampaignID:      r.CampaignID,
		CampaignName:    r.CampaignName,
		AdSetName:       r.AdSetName,
		Status:          r.Status,
		EffectiveStatus: r.EffectiveStatus,
		Budget:          resolveBudget(r.DailyBudget, r.LifetimeBudget),
		AmountSpent:     formatNullFloat(r.Spend),
		Impressions:     formatNullInt(r.Impressions),
		Reach:           formatNullInt(r.Reach),
	}

	if r.EndTime != nil {
		row.Ends = r.EndTime.Format("2006-01-02")
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

	if r.AttributionSpec != nil {
		var v interface{}
		if err := json.Unmarshal(r.AttributionSpec, &v); err == nil {
			row.AttributionSetting = v
		}
	}

	return row
}

type metaAction struct {
	ActionType string `json:"action_type"`
	Value      string `json:"value"`
}

func resolveBudget(daily, lifetime float64) string {
	if daily > 0 {
		return formatFloat(daily)
	}
	return formatFloat(lifetime)
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
