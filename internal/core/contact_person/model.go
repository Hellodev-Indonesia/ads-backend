package contact_person

import (
	"time"

	"gorm.io/gorm"
)

type ContactPerson struct {
	ID        uint64         `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:255;not null"`
	Phone     string         `json:"phone" gorm:"size:50;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (ContactPerson) TableName() string {
	return "contact_persons"
}
