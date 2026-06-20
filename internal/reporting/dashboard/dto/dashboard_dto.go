package dto

// SummaryResponse represents the response for the dashboard summary metrics.
type SummaryResponse struct {
	TotalSpend       float64 `json:"total_spend"`
	TotalOngoingAds  int64   `json:"total_ongoing_ads"`
	SecurityStatus   string  `json:"security_status"`
	SecurityIssues   int64   `json:"security_issues"`
}

// BrandMonitoringItem represents a single brand's monitoring metrics.
type BrandMonitoringItem struct {
	BrandID        uint64  `json:"brand_id"`
	BrandName      string  `json:"brand_name"`
	AdAccountCount int64   `json:"ad_account_count"`
	ActiveCampaign int64   `json:"active_campaign"`
	TotalSpends    float64 `json:"total_spends"`
	Photo          string  `json:"photo"`
}



// RecentActivityItem represents a single recent activity formatted for the dashboard.
type RecentActivityItem struct {
	ActivityID      uint64  `json:"activity_id"`
	BrandName       string  `json:"brand_name"`
	AdAccountName   string  `json:"ad_account_name"`
	Activity        string  `json:"activity"`
	ActivityDetails string  `json:"activity_details"`
	ItemChanged     string  `json:"item_changed"`
	ChangeBy        string  `json:"change_by"`
	DateAndTime     string  `json:"date_and_time"`
}


