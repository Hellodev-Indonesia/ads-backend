package ad_account

import (
	"time"
)

type MetaAdAccount struct {
	ID            string    `gorm:"column:id;primaryKey" json:"id"`
	Name          string    `gorm:"column:name;not null" json:"name"`
	AccountStatus int       `gorm:"column:account_status;not null;default:1" json:"account_status"`
	BrandID       *uint64   `gorm:"column:brand_id" json:"brand_id,omitempty"`
	Currency      *string   `gorm:"column:currency" json:"currency,omitempty"`
	TimezoneName  *string   `gorm:"column:timezone_name" json:"timezone_name,omitempty"`
	BusinessID    *string   `gorm:"column:business_id" json:"business_id,omitempty"`
	IsActive      bool      `gorm:"column:is_active;not null;default:true" json:"is_active"`
	SyncedAt      time.Time `gorm:"column:synced_at;autoUpdateTime" json:"synced_at"`
}

func (MetaAdAccount) TableName() string {
	return "meta_ad_accounts"
}
