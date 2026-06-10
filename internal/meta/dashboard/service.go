package dashboard

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/alex/ads_backend/internal/meta/dashboard/dto"
	"github.com/alex/ads_backend/pkg/response"
)

// Action type keys used by Meta Graph API
const (
	actionTotalMessaging = "onsite_conversion.total_messaging_connection"
	actionNewMessaging   = "onsite_conversion.messaging_first_reply"
	actionPurchase       = "purchase"
)

type metaAction struct {
	ActionType string `json:"action_type"`
	Value      string `json:"value"`
}

type Service interface {
	GetCampaignDashboard(filter DashboardFilter) ([]dto.CampaignDashboardRow, *response.PaginationMeta, error)
	GetAdSetDashboard(filter DashboardFilter) ([]dto.AdSetDashboardRow, *response.PaginationMeta, error)
	GetAdDashboard(filter DashboardFilter) ([]dto.AdDashboardRow, *response.PaginationMeta, error)
	GetBrandDashboard(filter DashboardFilter) ([]dto.BrandDashboardResponse, *response.PaginationMeta, error)
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo}
}

func (s *serviceImpl) GetCampaignDashboard(filter DashboardFilter) ([]dto.CampaignDashboardRow, *response.PaginationMeta, error) {
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

func formatNullInt(v *int64) string {
	if v == nil {
		return "0"
	}
	return strconv.FormatInt(*v, 10)
}

func formatFloat(v float64) string {
	if v == 0 {
		return "0"
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func (s *serviceImpl) GetAdSetDashboard(filter DashboardFilter) ([]dto.AdSetDashboardRow, *response.PaginationMeta, error) {
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

func (s *serviceImpl) GetAdDashboard(filter DashboardFilter) ([]dto.AdDashboardRow, *response.PaginationMeta, error) {
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

	return row
}

func (s *serviceImpl) GetBrandDashboard(filter DashboardFilter) ([]dto.BrandDashboardResponse, *response.PaginationMeta, error) {
	rows, total, err := s.repo.FindBrandDashboard(filter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch brand dashboard: %w", err)
	}

	result := make([]dto.BrandDashboardResponse, 0, len(rows))
	for _, r := range rows {
		var spend float64
		if r.TotalSpends != nil {
			spend = *r.TotalSpends
		}
		result = append(result, dto.BrandDashboardResponse{
			BrandID:             r.BrandID,
			BrandName:           r.BrandName,
			AdAccountCount:      r.AdAccountCount,
			ActiveCampaignCount: r.ActiveCampaignCount,
			TotalSpends:         spend,
		})
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
