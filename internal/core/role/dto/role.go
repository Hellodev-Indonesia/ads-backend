package dto

import "github.com/alex/ads_backend/internal/core/permission/dto"

type RoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type RoleResponse struct {
	ID          uint                     `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Permissions []dto.PermissionResponse `json:"permissions"`
}

type AssignPermissionRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}
