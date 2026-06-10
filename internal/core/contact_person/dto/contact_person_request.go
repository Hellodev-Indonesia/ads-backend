package dto

type ContactPersonRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=255" example:"John Doe"`
	Phone string `json:"phone" binding:"required,min=6,max=50" example:"+6281234567890"`
}
