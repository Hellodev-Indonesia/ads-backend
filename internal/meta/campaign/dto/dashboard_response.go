package dto

type CampaignDashboardRow struct {
	CampaignID      string `json:"campaign_id"`
	CampaignName    string `json:"campaign_name"`
	Status          string `json:"status"`
	EffectiveStatus string `json:"effective_status"`
	Objective       string `json:"objective,omitempty"`

	// Budget: whichever is set (daily takes precedence)
	Budget string `json:"budget"`

	// Insight metrics
	AmountSpent string `json:"amount_spent"`
	Impressions string `json:"impressions"`
	Reach       string `json:"reach"`

	// Derived from insights.actions JSON
	Results       string `json:"results"`
	CostPerResult string `json:"cost_per_result"`

	TotalMessagingConversations string `json:"total_messaging_conversations"`
	NewMessagingConnections     string `json:"new_messaging_connections"`
	Purchases                   string `json:"purchases"`

	// Campaign schedule
	Schedule string `json:"schedule,omitempty"`
	Ends     string `json:"ends,omitempty"`

	// From adset.attribution_spec
	AttributionSetting interface{} `json:"attribution_setting,omitempty"`

	BidStrategy         string `json:"bid_strategy,omitempty"`
	LastSignificantEdit string `json:"last_significant_edit,omitempty"`
	CostPerPurchase     string `json:"cost_per_purchase,omitempty"`

	DateStart string `json:"date_start"`
	DateStop  string `json:"date_stop"`
}
