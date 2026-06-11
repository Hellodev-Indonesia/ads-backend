package dto

type AdDashboardRow struct {
	AdID            string `json:"ad_id"`
	AdSetID         string `json:"adset_id"`
	CampaignID      string `json:"campaign_id"`
	CampaignName    string `json:"campaign_name"`
	AdSetName       string `json:"adset_name"`
	AdName          string `json:"ad_name"`
	Status          string `json:"status"`
	EffectiveStatus string `json:"effective_status"`
	CreativeID      string `json:"creative_id,omitempty"`

	AmountSpent string `json:"amount_spent"`
	Impressions string `json:"impressions"`
	Reach       string `json:"reach"`

	Results       string `json:"results"`
	CostPerResult string `json:"cost_per_result"`

	TotalMessagingConversations string `json:"total_messaging_conversations"`
	NewMessagingConnections     string `json:"new_messaging_connections"`
	Purchases                   string `json:"purchases"`

	DateStart string `json:"date_start"`
	DateStop  string `json:"date_stop"`
}
