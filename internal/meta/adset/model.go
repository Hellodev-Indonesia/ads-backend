package adset

import (
	"encoding/json"
	"time"
)

type MetaAdSet struct {
	ID              string          `gorm:"column:id;primaryKey" json:"id"`
	CampaignID      string          `gorm:"column:campaign_id;not null" json:"campaign_id"`
	Name            string          `gorm:"column:name;not null" json:"name"`
	Status          string          `gorm:"column:status;not null" json:"status"`
	EffectiveStatus string          `gorm:"column:effective_status" json:"effective_status"`
	DailyBudget     float64         `gorm:"column:daily_budget;default:0" json:"daily_budget"`
	LifetimeBudget  float64         `gorm:"column:lifetime_budget;default:0" json:"lifetime_budget"`
	BudgetRemaining float64         `gorm:"column:budget_remaining;default:0" json:"budget_remaining"`
	BidStrategy     string          `gorm:"column:bid_strategy" json:"bid_strategy"`
	AttributionSpec json.RawMessage `gorm:"column:attribution_spec;type:json" json:"attribution_spec"`
	StartTime       *time.Time      `gorm:"column:start_time" json:"start_time"`
	EndTime         *time.Time      `gorm:"column:end_time" json:"end_time"`
	CreatedTime     *time.Time      `gorm:"column:created_time" json:"created_time"`
	UpdatedTime     *time.Time      `gorm:"column:updated_time" json:"updated_time"`
	SyncedAt        time.Time       `gorm:"column:synced_at;autoUpdateTime" json:"synced_at"`
}

func (MetaAdSet) TableName() string {
	return "meta_ad_sets"
}
