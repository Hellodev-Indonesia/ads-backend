package dto

type CampaignSummaryResponse struct {
	AmountSpent     float64 `json:"amount_spent"`
	Impressions     int64   `json:"impressions"`
	Reach           int64   `json:"reach"`
	TotalMessaging  int64   `json:"total_messaging"`
	NewMessaging    int64   `json:"new_messaging"`
	PurchaseTotal   int64   `json:"purchase_total"`
	CostPerPurchase float64 `json:"cost_per_purchase"`
}
