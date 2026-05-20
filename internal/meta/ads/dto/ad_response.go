package dto

type CreativeRef struct {
	ID string `json:"id" example:"2387654321098"`
}

type AdResponse struct {
	ID              string      `json:"id" example:"2386543210987"`
	Name            string      `json:"name" example:"Promo Image Ad"`
	AdSetID         string      `json:"adset_id" example:"2389876543210"`
	CampaignID      string      `json:"campaign_id" example:"2381234567890"`
	Status          string      `json:"status" example:"ACTIVE"`
	EffectiveStatus string      `json:"effective_status" example:"ACTIVE"`
	Creative        CreativeRef `json:"creative"`
	CreatedTime     string      `json:"created_time,omitempty" example:"2026-05-11T10:00:00Z"`
	UpdatedTime     string      `json:"updated_time,omitempty" example:"2026-05-11T12:00:00Z"`
}
