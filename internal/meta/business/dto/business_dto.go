package dto

import "time"

type BusinessResponse struct {
	ID                string     `json:"id"`
	Name              string     `json:"name"`
	ProfilePictureURI *string    `json:"profile_picture_uri,omitempty"`
	TimezoneID        *int       `json:"timezone_id,omitempty"`
	CreatedTime       *time.Time `json:"created_time,omitempty"`
	SyncedAt          time.Time  `json:"synced_at"`
}
