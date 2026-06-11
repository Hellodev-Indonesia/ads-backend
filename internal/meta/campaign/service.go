package campaign

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/alex/ads_backend/internal/meta/campaign/dto"
	"github.com/alex/ads_backend/pkg/meta_client"
	"github.com/alex/ads_backend/pkg/response"
)

const DefaultFields = "id,name,status,effective_status,objective,daily_budget,lifetime_budget,budget_remaining,bid_strategy,start_time,stop_time,created_time,updated_time"

type Service interface {
	// DB reads (used by handlers)
	GetCampaigns(filter CampaignFilter) ([]dto.CampaignResponse, *response.PaginationMeta, error)
	GetCampaignByID(id string) (*dto.CampaignResponse, error)
	GetSummaryByBrand(brandID uint64, dateStart, dateStop string) (*dto.CampaignSummaryResponse, error)
	GetCampaignDashboard(filter CampaignFilter) ([]dto.CampaignDashboardRow, *response.PaginationMeta, error)

	// Meta API sync (used by sync job)
	SyncCampaigns(adAccountID string) (int, error)
	SyncCampaignsByIDs(adAccountID string, ids []string) (int, error)
}

type serviceImpl struct {
	client *meta_client.Client
	repo   Repository
}

func NewService(client *meta_client.Client, repo Repository) Service {
	return &serviceImpl{client: client, repo: repo}
}

// --- DB READ METHODS (for API handlers) ---

func (s *serviceImpl) GetCampaigns(filter CampaignFilter) ([]dto.CampaignResponse, *response.PaginationMeta, error) {
	campaigns, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch campaigns: %w", err)
	}

	var result []dto.CampaignResponse
	for _, c := range campaigns {
		result = append(result, mapModelToDTO(c))
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

func (s *serviceImpl) GetCampaignByID(id string) (*dto.CampaignResponse, error) {
	campaign, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch campaign: %w", err)
	}

	resp := mapModelToDTO(*campaign)
	return &resp, nil
}

type metaAction struct {
	ActionType string `json:"action_type"`
	Value      string `json:"value"`
}

func (s *serviceImpl) GetSummaryByBrand(brandID uint64, dateStart, dateStop string) (*dto.CampaignSummaryResponse, error) {
	rows, err := s.repo.GetSummaryByBrand(brandID, dateStart, dateStop)
	if err != nil {
		return nil, fmt.Errorf("failed to get summary: %w", err)
	}

	var summary dto.CampaignSummaryResponse
	for _, row := range rows {
		summary.AmountSpent += row.Spend
		summary.Impressions += row.Impressions
		summary.Reach += row.Reach

		if len(row.Actions) > 0 {
			var actions []metaAction
			if err := json.Unmarshal(row.Actions, &actions); err == nil {
				for _, act := range actions {
					val, _ := strconv.ParseInt(act.Value, 10, 64)
					switch act.ActionType {
					case "onsite_conversion.total_messaging_connection":
						summary.TotalMessaging += val
					case "onsite_conversion.messaging_first_reply":
						summary.NewMessaging += val
					case "purchase":
						summary.PurchaseTotal += val
					}
				}
			}
		}
	}

	if summary.PurchaseTotal > 0 {
		summary.CostPerPurchase = summary.AmountSpent / float64(summary.PurchaseTotal)
	}

	return &summary, nil
}

func (s *serviceImpl) GetCampaignDashboard(filter CampaignFilter) ([]dto.CampaignDashboardRow, *response.PaginationMeta, error) {
	rows, total, err := s.repo.FindCampaignDashboard(filter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch campaign dashboard: %w", err)
	}

	result := make([]dto.CampaignDashboardRow, 0, len(rows))
	for _, r := range rows {
		dtoRow := mapScanToDTO(r)
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

func mapScanToDTO(r campaignDashboardScan) dto.CampaignDashboardRow {
	var budgetStr string
	if r.DailyBudget == 0 && r.LifetimeBudget == 0 {
		budgetStr = resolveBudget(r.AdsetDailyBudget, r.AdsetLifetimeBudget)
	} else {
		budgetStr = resolveBudget(r.DailyBudget, r.LifetimeBudget)
	}

	row := dto.CampaignDashboardRow{
		CampaignID:      r.CampaignID,
		CampaignName:    r.CampaignName,
		Status:          r.Status,
		EffectiveStatus: r.EffectiveStatus,
		Objective:       r.Objective,
		Budget:          budgetStr,
		AmountSpent:     formatNullFloat(r.Spend),
		Impressions:     formatNullInt(r.Impressions),
		Reach:           formatNullInt(r.Reach),
		BidStrategy:     r.BidStrategy,
	}

	if r.UpdatedTime != nil {
		row.LastSignificantEdit = r.UpdatedTime.Format(time.RFC3339)
	}
	if r.StartTime != nil {
		row.Schedule = r.StartTime.Format("2006-01-02")
	}

	if r.StopTime != nil {
		row.Ends = r.StopTime.Format("2006-01-02")
	}
	if r.DateStart != nil && len(*r.DateStart) >= 10 {
		row.DateStart = (*r.DateStart)[:10]
	}
	if r.DateStop != nil && len(*r.DateStop) >= 10 {
		row.DateStop = (*r.DateStop)[:10]
	}

	// Extract action metrics from JSON
	actions := parseActions(r.Actions)
	row.TotalMessagingConversations = findAction(actions, "onsite_conversion.total_messaging_connection")
	row.NewMessagingConnections = findAction(actions, "onsite_conversion.messaging_first_reply")
	row.Purchases = findAction(actions, "purchase")

	var primaryActionType string
	if len(actions) > 0 {
		row.Results = actions[0].Value
		primaryActionType = actions[0].ActionType
	}

	if val := findAction(actions, "onsite_conversion.messaging_first_reply"); val != "0" {
		row.Results = val
		primaryActionType = "onsite_conversion.messaging_first_reply"
	} else if val := findAction(actions, "onsite_conversion.total_messaging_connection"); val != "0" {
		row.Results = val
		primaryActionType = "onsite_conversion.total_messaging_connection"
	} else if val := findAction(actions, "purchase"); val != "0" {
		row.Results = val
		primaryActionType = "purchase"
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

	// Calculate CostPerPurchase
	spent, _ := strconv.ParseFloat(row.AmountSpent, 64)
	purchases, _ := strconv.ParseFloat(row.Purchases, 64)
	if purchases > 0 {
		row.CostPerPurchase = formatFloat(math.Ceil(spent / purchases))
	} else {
		row.CostPerPurchase = "0"
	}

	// Attribution spec from adset
	if r.AttributionSpec != nil {
		var v interface{}
		if err := json.Unmarshal(r.AttributionSpec, &v); err == nil {
			row.AttributionSetting = v
		}
	}

	return row
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

// --- META API SYNC METHOD (for background job) ---

func (s *serviceImpl) SyncCampaigns(adAccountID string) (int, error) {
	params := url.Values{}
	params.Set("fields", DefaultFields)

	rawList, _, err := s.client.Get(adAccountID+"/campaigns", params, true)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch campaigns from Meta: %w", err)
	}

	var models []MetaCampaign
	for _, raw := range rawList {
		var item dto.CampaignResponse
		if err := json.Unmarshal(raw, &item); err != nil {
			log.Printf("Warning: skipping campaign unmarshal error: %v", err)
			continue
		}
		models = append(models, mapDTOToModel(item, adAccountID))
	}

	if err := s.repo.UpsertBatch(models); err != nil {
		return 0, fmt.Errorf("failed to upsert campaigns: %w", err)
	}

	return len(models), nil
}

func (s *serviceImpl) SyncCampaignsByIDs(adAccountID string, ids []string) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	batchSize := 50
	var allModels []MetaCampaign

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
			log.Printf("Warning: failed to bulk fetch campaigns: %v", err)
			continue
		}

		if len(rawList) == 0 {
			continue
		}

		var rawMap map[string]json.RawMessage
		if err := json.Unmarshal(rawList[0], &rawMap); err != nil {
			log.Printf("Warning: failed to unmarshal bulk campaigns response: %v", err)
			continue
		}

		for _, rawPayloadBytes := range rawMap {
			var item dto.CampaignResponse
			if err := json.Unmarshal(rawPayloadBytes, &item); err != nil {
				continue
			}
			allModels = append(allModels, mapDTOToModel(item, adAccountID))
		}
	}

	if len(allModels) > 0 {
		if err := s.repo.UpsertBatch(allModels); err != nil {
			return 0, fmt.Errorf("failed to upsert campaigns: %w", err)
		}
	}

	return len(allModels), nil
}

// --- MAPPERS ---

func mapModelToDTO(m MetaCampaign) dto.CampaignResponse {
	return dto.CampaignResponse{
		ID:              m.ID,
		AccountID:       m.AccountID,
		Name:            m.Name,
		Status:          m.Status,
		EffectiveStatus: m.EffectiveStatus,
		Objective:       m.Objective,
		BuyingType:      m.BuyingType,
		DailyBudget:     formatDecimal(m.DailyBudget),
		LifetimeBudget:  formatDecimal(m.LifetimeBudget),
		BudgetRemaining: formatDecimal(m.BudgetRemaining),
		SpendCap:        formatDecimal(m.SpendCap),
		BidStrategy:     m.BidStrategy,
		StartTime:       formatTime(m.StartTime),
		StopTime:        formatTime(m.StopTime),
		CreatedTime:     formatTime(m.CreatedTime),
		UpdatedTime:     formatTime(m.UpdatedTime),
	}
}

func mapDTOToModel(d dto.CampaignResponse, accountID string) MetaCampaign {
	acctID := d.AccountID
	if acctID == "" {
		acctID = accountID
	}
	return MetaCampaign{
		ID:              d.ID,
		AccountID:       acctID,
		Name:            d.Name,
		Status:          d.Status,
		EffectiveStatus: d.EffectiveStatus,
		Objective:       d.Objective,
		BuyingType:      d.BuyingType,
		DailyBudget:     parseDecimal(d.DailyBudget),
		LifetimeBudget:  parseDecimal(d.LifetimeBudget),
		BudgetRemaining: parseDecimal(d.BudgetRemaining),
		SpendCap:        parseDecimal(d.SpendCap),
		BidStrategy:     d.BidStrategy,
		StartTime:       parseTime(d.StartTime),
		StopTime:        parseTime(d.StopTime),
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
	// Try multiple formats that Meta uses
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
