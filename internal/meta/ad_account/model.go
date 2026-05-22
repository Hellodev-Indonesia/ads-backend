package ad_account

import (
	"time"
)

type MetaAdAccount struct {
	ID            string    `gorm:"column:id;primaryKey" json:"id"`
	Name          string    `gorm:"column:name;not null" json:"name"`
	AccountStatus int       `gorm:"column:account_status;not null;default:1" json:"account_status"`
	SyncedAt      time.Time `gorm:"column:synced_at;autoUpdateTime" json:"synced_at"`
}

func (MetaAdAccount) TableName() string {
	return "meta_ad_accounts"
}
