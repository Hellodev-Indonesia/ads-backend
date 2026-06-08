package dto

type AlertResponse struct {
	ID         uint64  `json:"id"`
	FraudLogID *uint64 `json:"fraud_log_id,omitempty"`
	BrandID    *uint64 `json:"brand_id,omitempty"`
	Title      string  `json:"title"`
	Message    string  `json:"message"`
	Severity   string  `json:"severity"`
	IsRead     bool    `json:"is_read"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type AlertFilter struct {
	Page      int     `form:"page"`
	Limit     int     `form:"limit"`
	BrandID   *uint64 `form:"brand_id"`
	Severity  string  `form:"severity"`
	IsRead    *bool   `form:"is_read"`
	DateStart string  `form:"date_start"`
	DateStop  string  `form:"date_stop"`
}

// CreateAlertInput is used internally (not from HTTP) to create an alert.
type CreateAlertInput struct {
	FraudLogID *uint64
	BrandID    *uint64
	Title      string
	Message    string
	Severity   string
}
