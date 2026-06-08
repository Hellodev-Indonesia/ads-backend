package dto

type InsightResponse struct {
	AccountID          string      `json:"account_id,omitempty" example:"541050504790549"`
	AccountName        string      `json:"account_name,omitempty" example:"My Ad Account"`
	AccountCurrency    string      `json:"account_currency,omitempty" example:"USD"`
	CampaignID         string      `json:"campaign_id" example:"2381234567890"`
	CampaignName       string      `json:"campaign_name" example:"Summer Sale Campaign"`
	AdSetID            string      `json:"adset_id,omitempty" example:"2389876543210"`
	AdSetName          string      `json:"adset_name,omitempty" example:"AdSet Leads USA"`
	AdID               string      `json:"ad_id,omitempty" example:"2386543210987"`
	AdName             string      `json:"ad_name,omitempty" example:"Promo Image Ad"`
	Objective          string      `json:"objective,omitempty" example:"OUTCOMES_LEADS"`
	Impressions        string      `json:"impressions" example:"15200"`
	Reach              string      `json:"reach,omitempty" example:"12000"`
	Clicks             string      `json:"clicks" example:"350"`
	InlineLinkClicks   string      `json:"inline_link_clicks,omitempty" example:"300"`
	InlineLinkClickCtr string      `json:"inline_link_click_ctr,omitempty" example:"2.0"`
	Spend              string      `json:"spend" example:"125.50"`
	CPC                string      `json:"cpc" example:"0.36"`
	CPM                string      `json:"cpm" example:"8.26"`
	CTR                string      `json:"ctr" example:"2.30"`
	Actions            interface{} `json:"actions,omitempty"`
	ActionValues       interface{} `json:"action_values,omitempty"`
	CostPerActionType  interface{} `json:"cost_per_action_type,omitempty"`
	DateStart          string      `json:"date_start,omitempty" example:"2026-04-11"`
	DateStop           string      `json:"date_stop,omitempty" example:"2026-05-11"`
}
