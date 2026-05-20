package dto

type AdSetResponse struct {
	ID              string `json:"id" example:"2389876543210"`
	Name            string `json:"name" example:"AdSet Leads USA"`
	CampaignID      string `json:"campaign_id" example:"2381234567890"`
	Status          string `json:"status" example:"ACTIVE"`
	EffectiveStatus string `json:"effective_status" example:"ACTIVE"`
	DailyBudget     string `json:"daily_budget" example:"50000"`
	LifetimeBudget  string `json:"lifetime_budget" example:"0"`
	StartTime       string `json:"start_time" example:"2026-05-11T10:00:00Z"`
	EndTime         string `json:"end_time" example:"2026-06-11T10:00:00Z"`
}
