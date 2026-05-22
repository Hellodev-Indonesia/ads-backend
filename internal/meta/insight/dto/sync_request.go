package dto

type SyncInsightRequest struct {
	AdAccountID   string `json:"ad_account_id,omitempty" example:"act_123"`
	Level         string `json:"level,omitempty" example:"campaign"`
	DatePreset    string `json:"date_preset,omitempty" example:"last_30d"`
	DateStart     string `json:"date_start,omitempty" example:"2026-05-15"`
	DateStop      string `json:"date_stop,omitempty" example:"2026-05-15"`
	TimeIncrement int    `json:"time_increment,omitempty" example:"1"`
	Force         bool   `json:"force,omitempty" example:"false"`
}
