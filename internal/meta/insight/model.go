package insight

import (
	"encoding/json"
	"time"
)

type MetaInsight struct {
	ID                 uint            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	AccountID          string          `gorm:"column:account_id;not null" json:"account_id"`
	AccountName        string          `gorm:"column:account_name" json:"account_name"`
	AccountCurrency    string          `gorm:"column:account_currency" json:"account_currency"`
	CampaignID         string          `gorm:"column:campaign_id;not null" json:"campaign_id"`
	CampaignName       string          `gorm:"column:campaign_name" json:"campaign_name"`
	AdSetID            string          `gorm:"column:adset_id;not null;default:''" json:"adset_id"`
	AdSetName          string          `gorm:"column:adset_name" json:"adset_name"`
	AdID               string          `gorm:"column:ad_id;not null;default:''" json:"ad_id"`
	AdName             string          `gorm:"column:ad_name" json:"ad_name"`
	Level              string          `gorm:"column:level;not null" json:"level"`
	Objective          string          `gorm:"column:objective" json:"objective"`
	Impressions        int64           `gorm:"column:impressions;default:0" json:"impressions"`
	Reach              int64           `gorm:"column:reach;default:0" json:"reach"`
	Clicks             int64           `gorm:"column:clicks;default:0" json:"clicks"`
	InlineLinkClicks   int64           `gorm:"column:inline_link_clicks;default:0" json:"inline_link_clicks"`
	InlineLinkClickCtr float64         `gorm:"column:inline_link_click_ctr;default:0" json:"inline_link_click_ctr"`
	Spend              float64         `gorm:"column:spend;default:0" json:"spend"`
	CPC                float64         `gorm:"column:cpc;default:0" json:"cpc"`
	CPM                float64         `gorm:"column:cpm;default:0" json:"cpm"`
	CTR                float64         `gorm:"column:ctr;default:0" json:"ctr"`
	Actions            json.RawMessage `gorm:"column:actions;type:json" json:"actions"`
	ActionValues       json.RawMessage `gorm:"column:action_values;type:json" json:"action_values"`
	CostPerActionType  json.RawMessage `gorm:"column:cost_per_action_type;type:json" json:"cost_per_action_type"`
	DateStart          string          `gorm:"column:date_start;not null" json:"date_start"`
	DateStop           string          `gorm:"column:date_stop;not null" json:"date_stop"`
	SyncedAt           time.Time       `gorm:"column:synced_at;autoUpdateTime" json:"synced_at"`
}

func (MetaInsight) TableName() string {
	return "meta_insights"
}
