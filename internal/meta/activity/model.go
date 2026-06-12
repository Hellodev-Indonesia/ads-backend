package activity

import (
	"encoding/json"
	"time"
)

type MetaActivity struct {
	ID          string          `gorm:"primaryKey;column:id"`
	AdAccountID string          `gorm:"column:ad_account_id"`
	ActorID     *string         `gorm:"column:actor_id"`
	ActorName   *string         `gorm:"column:actor_name"`
	ObjectID    *string         `gorm:"column:object_id"`
	ObjectName  *string         `gorm:"column:object_name"`
	ObjectType  *string         `gorm:"column:object_type"`
	EventType   *string         `gorm:"column:event_type"`
	EventTime   *time.Time      `gorm:"column:event_time"`
	ExtraData   json.RawMessage `gorm:"column:extra_data;type:json"`
	CreatedAt   time.Time       `gorm:"column:created_at"`
	UpdatedAt   time.Time       `gorm:"column:updated_at"`
}

func (MetaActivity) TableName() string {
	return "meta_activities"
}
