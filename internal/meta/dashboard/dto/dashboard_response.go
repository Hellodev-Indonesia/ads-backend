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
	Ends string `json:"ends,omitempty"`

	// From adset.attribution_spec
	AttributionSetting interface{} `json:"attribution_setting,omitempty"`

	DateStart string `json:"date_start,omitempty"`
	DateStop  string `json:"date_stop,omitempty"`
}
