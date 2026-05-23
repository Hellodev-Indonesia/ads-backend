package ad_creative

import (
	"time"

	"gorm.io/gorm"
)

type AdCreative struct {
	ID               uint64         `gorm:"primaryKey" json:"id"`
	CreativeID       string         `gorm:"size:100;not null;uniqueIndex" json:"creative_id"`
	Name             *string        `gorm:"size:255" json:"name"`
	Title            *string        `gorm:"type:text" json:"title"`
	Body             *string        `gorm:"type:text" json:"body"`
	ImageURL         *string        `gorm:"type:text" json:"image_url"`
	VideoURL         *string        `gorm:"type:text" json:"video_url"`
	DestinationURL   *string        `gorm:"type:text" json:"destination_url"`
	NormalizedDomain *string        `gorm:"size:255;index" json:"normalized_domain"`
	URLHash          *string        `gorm:"size:64;index" json:"url_hash"`
	RawPayload       *string        `gorm:"type:json" json:"raw_payload"`
	SyncedAt         *time.Time     `json:"synced_at"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AdCreative) TableName() string {
	return "ad_creatives"
}

type AdCreativeVersion struct {
	ID               uint64     `gorm:"primaryKey" json:"id"`
	CreativeID       string     `gorm:"size:100;not null;index" json:"creative_id"`
	AdID             *string    `gorm:"size:100;index" json:"ad_id"`
	AdsetID          *string    `gorm:"size:100" json:"adset_id"`
	CampaignID       *string    `gorm:"size:100;index" json:"campaign_id"`
	AdAccountID      *string    `gorm:"size:100;index" json:"ad_account_id"`
	BrandID          *uint64    `gorm:"index" json:"brand_id"`
	Name             *string    `gorm:"size:255" json:"name"`
	Title            *string    `gorm:"type:text" json:"title"`
	Body             *string    `gorm:"type:text" json:"body"`
	ImageURL         *string    `gorm:"type:text" json:"image_url"`
	VideoURL         *string    `gorm:"type:text" json:"video_url"`
	DestinationURL   *string    `gorm:"type:text" json:"destination_url"`
	NormalizedDomain *string    `gorm:"size:255" json:"normalized_domain"`
	URLHash          *string    `gorm:"size:64;index" json:"url_hash"`
	RawPayload       *string    `gorm:"type:json" json:"raw_payload"`
	ChangedFields    *string    `gorm:"type:json" json:"changed_fields"`
	ChangeType       *string    `gorm:"size:100" json:"change_type"`
	SyncedAt         *time.Time `json:"synced_at"`
	CreatedAt        time.Time  `gorm:"index" json:"created_at"`
}

func (AdCreativeVersion) TableName() string {
	return "ad_creative_versions"
}
