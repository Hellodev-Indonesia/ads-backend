package dto

type InsightResponse struct {
	CampaignID   string `json:"campaign_id" example:"2381234567890"`
	CampaignName string `json:"campaign_name" example:"Summer Sale Campaign"`
	Impressions  string `json:"impressions" example:"15200"`
	Clicks       string `json:"clicks" example:"350"`
	Spend        string `json:"spend" example:"125.50"`
	CPC          string `json:"cpc" example:"0.36"`
	CPM          string `json:"cpm" example:"8.26"`
	CTR          string `json:"ctr" example:"2.30"`
}
