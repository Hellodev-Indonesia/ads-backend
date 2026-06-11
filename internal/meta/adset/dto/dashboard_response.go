package dto

type AdSetDashboardRow struct {
	AdSetID         string `json:"adset_id"`
	CampaignID      string `json:"campaign_id"`
	CampaignName    string `json:"campaign_name"`
	AdSetName       string `json:"adset_name"`
	Status          string `json:"status"`
	EffectiveStatus string `json:"effective_status"`

	Budget string `json:"budget"`

	AmountSpent string `json:"amount_spent"`
	Impressions string `json:"impressions"`
	Reach       string `json:"reach"`

	Results       string `json:"results"`
	CostPerResult string `json:"cost_per_result"`

	TotalMessagingConversations string `json:"total_messaging_conversations"`
	NewMessagingConnections     string `json:"new_messaging_connections"`
	Purchases                   string `json:"purchases"`

	Schedule string `json:"schedule,omitempty"`
	Ends     string `json:"ends,omitempty"`

	AttributionSetting interface{} `json:"attribution_setting,omitempty"`

	BidStrategy         string `json:"bid_strategy,omitempty"`
	LastSignificantEdit string `json:"last_significant_edit,omitempty"`
	CostPerPurchase     string `json:"cost_per_purchase,omitempty"`

	DateStart string `json:"date_start"`
	DateStop  string `json:"date_stop"`
}
