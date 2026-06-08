package dto

type TriggerSyncRequest struct {
	AdAccountID string `json:"ad_account_id,omitempty" example:"act_1234567890"`
	DateStart   string `json:"date_start,omitempty" example:"2026-03-01"`
	DateStop    string `json:"date_stop,omitempty" example:"2026-03-31"`
}
