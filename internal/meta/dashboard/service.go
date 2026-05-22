package dashboard

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/alex/ads_backend/internal/meta/dashboard/dto"
	"github.com/alex/ads_backend/pkg/response"
)

// Action type keys used by Meta Graph API
const (
	actionTotalMessaging = "onsite_conversion.total_messaging_connection"
	actionNewMessaging   = "onsite_conversion.new_messaging_connection"
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
	row := dto.CampaignDashboardRow{
		CampaignID:      r.CampaignID,
		CampaignName:    r.CampaignName,
		Status:          r.Status,
		EffectiveStatus: r.EffectiveStatus,
		Objective:       r.Objective,
		Budget:          resolveBudget(r.DailyBudget, r.LifetimeBudget),
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

	// Results = first action value (primary objective metric)
	if len(actions) > 0 {
		row.Results = actions[0].Value
	}

	// Cost per result = first cost_per_action_type value
	costs := parseActions(r.CostPerActionType)
	if len(costs) > 0 {
		row.CostPerResult = costs[0].Value
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

	if len(actions) > 0 {
		row.Results = actions[0].Value
	}

	costs := parseActions(r.CostPerActionType)
	if len(costs) > 0 {
		row.CostPerResult = costs[0].Value
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

	if len(actions) > 0 {
		row.Results = actions[0].Value
	}

	costs := parseActions(r.CostPerActionType)
	if len(costs) > 0 {
		row.CostPerResult = costs[0].Value
	}

	return row
}
