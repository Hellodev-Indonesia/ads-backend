package dto

type AdSetResponse struct {
	ID               string `json:"id" example:"2389876543210"`
	Name             string `json:"name" example:"AdSet Leads USA"`
	CampaignID       string `json:"campaign_id" example:"2381234567890"`
	Status           string `json:"status" example:"ACTIVE"`
	EffectiveStatus  string `json:"effective_status" example:"ACTIVE"`
	DailyBudget      string `json:"daily_budget,omitempty" example:"50000"`
	LifetimeBudget   string `json:"lifetime_budget,omitempty" example:"0"`
	BudgetRemaining  string `json:"budget_remaining,omitempty" example:"50000"`
	BidStrategy      string      `json:"bid_strategy,omitempty" example:"LOWEST_COST_WITHOUT_CAP"`
	AttributionSpec  interface{} `json:"attribution_spec,omitempty"`
	StartTime        string `json:"start_time,omitempty" example:"2026-05-11T10:00:00Z"`
	EndTime          string `json:"end_time,omitempty" example:"2026-06-11T10:00:00Z"`
	CreatedTime      string `json:"created_time,omitempty" example:"2026-05-11T10:00:00Z"`
	UpdatedTime      string `json:"updated_time,omitempty" example:"2026-05-11T12:00:00Z"`
}
