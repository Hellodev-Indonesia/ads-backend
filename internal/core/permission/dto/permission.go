package dto

type PermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type PermissionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PermissionFilter struct {
	Name  string `form:"name"`
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
}
