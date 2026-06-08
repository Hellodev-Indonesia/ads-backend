package campaign

import (
	"time"
)

type MetaCampaign struct {
	ID              string     `gorm:"column:id;primaryKey" json:"id"`
	AccountID       string     `gorm:"column:account_id;not null" json:"account_id"`
	Name            string     `gorm:"column:name;not null" json:"name"`
	Status          string     `gorm:"column:status;not null" json:"status"`
	EffectiveStatus string     `gorm:"column:effective_status" json:"effective_status"`
	Objective       string     `gorm:"column:objective" json:"objective"`
	BuyingType      string     `gorm:"column:buying_type" json:"buying_type"`
	DailyBudget     float64    `gorm:"column:daily_budget;default:0" json:"daily_budget"`
	LifetimeBudget  float64    `gorm:"column:lifetime_budget;default:0" json:"lifetime_budget"`
	BudgetRemaining float64    `gorm:"column:budget_remaining;default:0" json:"budget_remaining"`
	SpendCap        float64    `gorm:"column:spend_cap;default:0" json:"spend_cap"`
	BidStrategy     string     `gorm:"column:bid_strategy" json:"bid_strategy"`
	StartTime       *time.Time `gorm:"column:start_time" json:"start_time"`
	StopTime        *time.Time `gorm:"column:stop_time" json:"stop_time"`
	CreatedTime     *time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedTime     *time.Time `gorm:"column:updated_time" json:"updated_time"`
	SyncedAt        time.Time  `gorm:"column:synced_at;autoUpdateTime" json:"synced_at"`
}

func (MetaCampaign) TableName() string {
	return "meta_campaigns"
}
