package dto

type AssignBrandRequest struct {
	BrandID uint64 `json:"brand_id" binding:"required"`
}
