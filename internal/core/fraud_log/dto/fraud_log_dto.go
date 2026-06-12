package dto

type FraudLogResponse struct {
	ID            uint64  `json:"id"`
	Brand         *SimpleBrand     `json:"brand,omitempty"`
	AdAccount     *SimpleAdAccount `json:"ad_account,omitempty"`
	Campaign      *SimpleCampaign `json:"campaign,omitempty"`
	Adset         *SimpleAdSet    `json:"adset,omitempty"`
	Ad            *SimpleAd       `json:"ad,omitempty"`
	CreativeID    *string         `json:"creative_id,omitempty"`
	EventType     string  `json:"event_type"`
	Severity      string  `json:"severity"`
	OldValue      *string `json:"old_value,omitempty"`
	NewValue      *string `json:"new_value,omitempty"`
	MatchedRuleID *uint64 `json:"matched_rule_id,omitempty"`
	Message       *string `json:"message,omitempty"`
	Status        string  `json:"status"`
	DetectedAt    *string `json:"detected_at,omitempty"`
	ResolvedAt    *string `json:"resolved_at,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type SimpleBrand struct {
	ID    uint64  `json:"id"`
	Name  string  `json:"name"`
	Photo *string `json:"photo"`
}

type SimpleAdAccount struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	BusinessName *string `json:"business_name,omitempty"`
}

type SimpleCampaign struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SimpleAdSet struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SimpleAd struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type FraudLogFilter struct {
	Page        int     `form:"page"`
	Limit       int     `form:"limit"`
	BrandID     *uint64 `form:"brand_id"`
	AdAccountID string  `form:"ad_account_id"`
	CampaignID  string  `form:"campaign_id"`
	CreativeID  string  `form:"creative_id"`
	Severity    string  `form:"severity"`
	Status      string  `form:"status"`
	DateStart   string  `form:"date_start"`
	DateStop    string  `form:"date_stop"`
}

// CreateFraudLogInput is used internally (not from HTTP) to create a fraud log.
type CreateFraudLogInput struct {
	BrandID       *uint64
	AdAccountID   *string
	CampaignID    *string
	AdsetID       *string
	AdID          *string
	CreativeID    *string
	EventType     string
	Severity      string
	OldValue      *string
	NewValue      *string
	MatchedRuleID *uint64
	Message       string
}
