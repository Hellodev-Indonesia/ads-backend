package dto

type CampaignResponse struct {
	ID              string `json:"id" example:"2381234567890"`
	Name            string `json:"name" example:"Summer Sale Campaign"`
	Status          string `json:"status" example:"ACTIVE"`
	EffectiveStatus string `json:"effective_status" example:"ACTIVE"`
	Objective       string `json:"objective" example:"OUTCOMES_LEADS"`
	CreatedTime     string `json:"created_time" example:"2026-05-11T10:00:00Z"`
	UpdatedTime     string `json:"updated_time" example:"2026-05-11T12:00:00Z"`
}
