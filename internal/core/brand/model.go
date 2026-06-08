package brand

import (
	"time"

	"gorm.io/gorm"
)

type Brand struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Photo       *string        `gorm:"size:255" json:"photo,omitempty"`
	Description *string        `type:"text" json:"description,omitempty"`
	IsActive    bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	AdAccountCount int64          `gorm:"->;-:migration" json:"ad_account_count"`
}

func (Brand) TableName() string {
	return "brands"
}
