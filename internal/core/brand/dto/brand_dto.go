package dto

import "mime/multipart"

type CreateBrandRequest struct {
	Name        string                `form:"name" binding:"required,max=255"`
	Photo       *multipart.FileHeader `form:"photo" binding:"omitempty"`
	Description *string               `form:"description" binding:"omitempty"`
	IsActive    *bool                 `form:"is_active" binding:"omitempty"`
}

type UpdateBrandRequest struct {
	Name        *string               `form:"name" binding:"omitempty,max=255"`
	Photo       *multipart.FileHeader `form:"photo" binding:"omitempty"`
	Description *string               `form:"description" binding:"omitempty"`
	IsActive    *bool                 `form:"is_active" binding:"omitempty"`
}

type BrandResponse struct {
	ID             uint64  `json:"id"`
	Slug           string  `json:"slug"`
	Name           string  `json:"name"`
	Photo          *string `json:"photo,omitempty"`
	Description    *string `json:"description,omitempty"`
	IsActive       bool    `json:"is_active"`
	AdAccountCount int64   `json:"ad_account_count"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type BrandListResponse struct {
	Data []BrandResponse `json:"data"`
}

type BrandFilter struct {
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	Name     string `form:"name"`
	IsActive *bool  `form:"is_active"`
}
