package dto

type AssignBrandRequest struct {
	BrandID      *uint64  `json:"brand_id"`
	BusinessID   *string  `json:"business_id,omitempty"`
	AdAccountIDs []string `json:"ad_account_ids,omitempty"`
}
