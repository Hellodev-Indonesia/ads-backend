package dto

type AdsetSummaryResponse struct {
	AmountSpent     float64 `json:"amount_spent"`
	Impressions     int64   `json:"impressions"`
	Reach           int64   `json:"reach"`
	CostPerResult   float64 `json:"cost_per_result"`
	CostPerPurchase float64 `json:"cost_per_purchase"`
	TotalMessaging  int64   `json:"total_messaging"`
	NewMessaging    int64   `json:"new_messaging"`
	PurchaseTotal   int64   `json:"purchase_total"`
}
