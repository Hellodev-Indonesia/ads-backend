package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthUserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Token           string           `json:"token"`
	CentrifugoToken string           `json:"centrifugo_token"`
	User            AuthUserResponse `json:"user"`
	Roles           []string         `json:"roles"`
	Permissions     []string         `json:"permissions"`
}
