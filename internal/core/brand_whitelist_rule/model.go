package brand_whitelist_rule

import (
	"time"

	"gorm.io/gorm"
)

type BrandWhitelistRule struct {
	ID               uint64         `gorm:"primaryKey" json:"id"`
	BrandID          uint64         `gorm:"index;not null" json:"brand_id"`
	Scope            string         `gorm:"size:50;not null;index" json:"scope"`
	MatchType        string         `gorm:"size:50;not null;index" json:"match_type"`
	Value            string         `gorm:"type:text;not null" json:"value"`
	NormalizedValue  *string        `gorm:"type:text" json:"normalized_value"`
	AllowSubdomains  bool           `gorm:"not null;default:false" json:"allow_subdomains"`
	IsActive         bool           `gorm:"not null;default:true;index" json:"is_active"`
	Description      *string        `gorm:"type:text" json:"description"`
	CreatedBy        *uint64        `json:"created_by"`
	ApprovedBy       *uint64        `json:"approved_by"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BrandWhitelistRule) TableName() string {
	return "brand_whitelist_rules"
}
