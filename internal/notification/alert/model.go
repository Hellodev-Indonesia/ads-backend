package alert

import (
	"time"
)

type Alert struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	FraudLogID *uint64   `json:"fraud_log_id"`
	BrandID    *uint64   `gorm:"index" json:"brand_id"`
	Title      string    `gorm:"size:255;not null" json:"title"`
	Message    string    `gorm:"type:text;not null" json:"message"`
	Severity   string    `gorm:"size:50;not null" json:"severity"`
	IsRead     bool      `gorm:"not null;default:false;index" json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Alert) TableName() string {
	return "alerts"
}
