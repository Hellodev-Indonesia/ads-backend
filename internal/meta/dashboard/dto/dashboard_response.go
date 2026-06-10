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

	DateStart string `json:"date_start"`
	DateStop  string `json:"date_stop"`
}

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

	Ends string `json:"ends,omitempty"`

	AttributionSetting interface{} `json:"attribution_setting,omitempty"`

	DateStart string `json:"date_start"`
	DateStop  string `json:"date_stop"`
}

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

type BrandDashboardResponse struct {
	BrandID             uint64  `json:"brand_id"`
	BrandName           string  `json:"brand_name"`
	AdAccountCount      int     `json:"ad_account_count"`
	ActiveCampaignCount int     `json:"active_campaign_count"`
	TotalSpends         float64 `json:"total_spends"`
}
