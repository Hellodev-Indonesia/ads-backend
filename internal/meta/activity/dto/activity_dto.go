package dto

type ActivityResponse struct {
	AdAccount       string `json:"ad_account"`
	Activity        string `json:"activity"`
	ActivityDetails string `json:"activity_details"`
	ItemChanged     string `json:"item_changed"`
	ChangeBy        string `json:"change_by"`
	DateAndTime     string `json:"date_and_time"`
}

type ActivityFilter struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}
