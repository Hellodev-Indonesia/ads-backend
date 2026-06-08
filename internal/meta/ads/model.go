package ads

import (
	"time"
)

type MetaAd struct {
	ID              string     `gorm:"column:id;primaryKey" json:"id"`
	CampaignID      string     `gorm:"column:campaign_id;not null" json:"campaign_id"`
	AdSetID         string     `gorm:"column:adset_id;not null" json:"adset_id"`
	Name            string     `gorm:"column:name;not null" json:"name"`
	Status          string     `gorm:"column:status;not null" json:"status"`
	EffectiveStatus string     `gorm:"column:effective_status" json:"effective_status"`
	CreativeID      string     `gorm:"column:creative_id" json:"creative_id"`
	CreatedTime     *time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedTime     *time.Time `gorm:"column:updated_time" json:"updated_time"`
	SyncedAt        time.Time  `gorm:"column:synced_at;autoUpdateTime" json:"synced_at"`
}

func (MetaAd) TableName() string {
	return "meta_ads"
}
