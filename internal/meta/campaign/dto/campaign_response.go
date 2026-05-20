package dto

type CampaignResponse struct {
	ID              string `json:"id" example:"2381234567890"`
	AccountID       string `json:"account_id,omitempty" example:"541050504790549"`
	Name            string `json:"name" example:"Summer Sale Campaign"`
	Status          string `json:"status" example:"ACTIVE"`
	EffectiveStatus string `json:"effective_status,omitempty" example:"ACTIVE"`
	Objective       string `json:"objective" example:"OUTCOMES_LEADS"`
	BuyingType      string `json:"buying_type,omitempty" example:"AUCTION"`
	DailyBudget     string `json:"daily_budget,omitempty" example:"50000"`
	LifetimeBudget  string `json:"lifetime_budget,omitempty" example:"0"`
	BudgetRemaining string `json:"budget_remaining,omitempty" example:"50000"`
	SpendCap        string `json:"spend_cap,omitempty" example:"0"`
	BidStrategy     string `json:"bid_strategy,omitempty" example:"LOWEST_COST_WITHOUT_CAP"`
	StartTime       string `json:"start_time,omitempty" example:"2026-05-11T10:00:00Z"`
	StopTime        string `json:"stop_time,omitempty" example:"2026-06-11T10:00:00Z"`
	CreatedTime     string `json:"created_time,omitempty" example:"2026-05-11T10:00:00Z"`
	UpdatedTime     string `json:"updated_time,omitempty" example:"2026-05-11T12:00:00Z"`
}
