package dto

type CreateWhitelistRuleRequest struct {
	Scope           string  `json:"scope" binding:"required,oneof=destination_url display_url url_tags domain"`
	MatchType       string  `json:"match_type" binding:"required,oneof=exact_url domain path_prefix contains regex"`
	Value           string  `json:"value" binding:"required"`
	AllowSubdomains bool    `json:"allow_subdomains"`
	IsActive        *bool   `json:"is_active"`
	Description     *string `json:"description"`
}

type UpdateWhitelistRuleRequest struct {
	Scope           *string `json:"scope" binding:"omitempty,oneof=destination_url display_url url_tags domain"`
	MatchType       *string `json:"match_type" binding:"omitempty,oneof=exact_url domain path_prefix contains regex"`
	Value           *string `json:"value"`
	AllowSubdomains *bool   `json:"allow_subdomains"`
	IsActive        *bool   `json:"is_active"`
	Description     *string `json:"description"`
}

type WhitelistRuleResponse struct {
	ID              uint64  `json:"id"`
	BrandID         uint64  `json:"brand_id"`
	Scope           string  `json:"scope"`
	MatchType       string  `json:"match_type"`
	Value           string  `json:"value"`
	NormalizedValue *string `json:"normalized_value,omitempty"`
	AllowSubdomains bool    `json:"allow_subdomains"`
	IsActive        bool    `json:"is_active"`
	Description     *string `json:"description,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type WhitelistRuleFilter struct {
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
	Scope     string `form:"scope"`
	MatchType string `form:"match_type"`
	IsActive  *bool  `form:"is_active"`
}

type CheckURLRequest struct {
	URL   string `json:"url" binding:"required"`
	Scope string `json:"scope" binding:"required,oneof=destination_url display_url url_tags domain"`
}

type CheckURLResponse struct {
	Allowed     bool                   `json:"allowed"`
	MatchedRule *WhitelistRuleResponse `json:"matched_rule,omitempty"`
}
