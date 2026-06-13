package fraud_log

import (
	"time"
)

type FraudLog struct {
	ID            uint64     `gorm:"primaryKey" json:"id"`
	BrandID       *uint64    `gorm:"index" json:"brand_id"`
	AdAccountID   *string    `gorm:"size:100;index" json:"ad_account_id"`
	CampaignID    *string    `gorm:"size:100;index" json:"campaign_id"`
	AdsetID       *string    `gorm:"size:100" json:"adset_id"`
	AdID          *string    `gorm:"size:100" json:"ad_id"`
	CreativeID    *string    `gorm:"size:100;index" json:"creative_id"`
	EventType     string     `gorm:"size:100;not null" json:"event_type"`
	ActorID       *string    `gorm:"size:255" json:"actor_id"`
	ActorName     *string    `gorm:"size:255" json:"actor_name"`
	Severity      string     `gorm:"size:50;not null" json:"severity"`
	OldValue      *string    `gorm:"type:text" json:"old_value"`
	NewValue      *string    `gorm:"type:text" json:"new_value"`
	MatchedRuleID *uint64    `json:"matched_rule_id"`
	Message       *string    `gorm:"type:text" json:"message"`
	Status        string     `gorm:"size:50;not null;default:'open';index" json:"status"`
	DetectedAt    *time.Time `json:"detected_at"`
	ResolvedAt    *time.Time `json:"resolved_at"`
	ResolvedBy    *uint64    `json:"resolved_by"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (FraudLog) TableName() string {
	return "fraud_logs"
}
