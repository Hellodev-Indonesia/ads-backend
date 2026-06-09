package dto

type AdAccountResponse struct {
	ID            string  `json:"id" example:"act_123456789"`
	Name          string  `json:"name" example:"My Ad Account"`
	AccountStatus int     `json:"account_status" example:"1"`
	BrandID       *uint64 `json:"brand_id,omitempty"`
	Currency      *string `json:"currency,omitempty"`
	TimezoneName  *string `json:"timezone_name,omitempty"`
	BusinessID    *string `json:"business_id,omitempty"`
	BusinessName  *string `json:"business_name,omitempty"`
	IsActive      bool    `json:"is_active"`
}
