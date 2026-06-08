package dto

import "time"

type UserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
	RoleIDs  []uint `json:"role_ids"`
}

type UserFilter struct {
	Name  string `form:"name"`
	Email string `form:"email"`
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
}

type UserResponse struct {
	ID        uint        `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Roles     []RoleBrief `json:"roles,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

type RoleBrief struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
