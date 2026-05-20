package dto

type AdAccountResponse struct {
	ID            string `json:"id" example:"act_123456789"`
	Name          string `json:"name" example:"My Ad Account"`
	AccountStatus int    `json:"account_status" example:"1"`
}
