package dto

type BrandDashboardResponse struct {
	Brand               BrandDashboardInfo `json:"brand"`
	AdAccountCount      int                `json:"ad_account_count"`
	ActiveCampaignCount int                `json:"active_campaign_count"`
	TotalSpends         float64            `json:"total_spends"`
}

type BrandDashboardInfo struct {
	ID    uint64  `json:"id"`
	Name  string  `json:"name"`
	Slug  string  `json:"slug"`
	Photo *string `json:"photo"`
}
