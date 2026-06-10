package dto

type TriggerSyncRequest struct {
	AdAccountID string `json:"ad_account_id,omitempty" form:"ad_account_id" example:"act_1234567890"`
	DateStart   string `json:"date_start,omitempty" form:"date_start" example:"2026-03-01"`
	DateStop    string `json:"date_stop,omitempty" form:"date_stop" example:"2026-03-31"`
}
