package dto

import "time"

type ContactPersonResponse struct {
	ID        uint64    `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Phone     string    `json:"phone" example:"+6281234567890"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ContactPersonListResponse struct {
	ID        uint64    `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Phone     string    `json:"phone" example:"+6281234567890"`
	CreatedAt time.Time `json:"created_at"`
}
