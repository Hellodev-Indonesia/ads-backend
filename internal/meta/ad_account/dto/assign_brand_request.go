package dto

type AssignBrandRequest struct {
	BrandID      *uint64  `json:"brand_id"`
	AdAccountIDs []string `json:"ad_account_ids,omitempty"`
}
