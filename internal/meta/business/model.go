package business

import "time"

type MetaBusiness struct {
	ID                string     `gorm:"primaryKey;size:50" json:"id"`
	Name              string     `gorm:"size:255;not null" json:"name"`
	ProfilePictureURI *string    `gorm:"type:text" json:"profile_picture_uri"`
	TimezoneID        *int       `json:"timezone_id"`
	CreatedTime       *time.Time `json:"created_time"`
	SyncedAt          time.Time  `gorm:"autoUpdateTime" json:"synced_at"`
}

func (MetaBusiness) TableName() string {
	return "meta_businesses"
}
