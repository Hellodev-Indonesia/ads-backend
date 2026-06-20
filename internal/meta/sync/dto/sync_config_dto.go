package dto

type MetaSyncConfigResponse struct {
	IntervalMinutes int  `json:"interval_minutes"`
	IsActive        bool `json:"is_active"`
}

type UpdateMetaSyncConfigRequest struct {
	IntervalMinutes int  `json:"interval_minutes" binding:"required,min=5"`
	IsActive        bool `json:"is_active"`
}
