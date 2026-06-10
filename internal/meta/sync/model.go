package sync

import (
	"errors"
	"time"

	"gorm.io/datatypes"
)

var ErrAlreadyRunning = errors.New("sync is already in progress")

const Channel = "meta:sync"

const (
	StatusPending       = "PENDING"
	StatusRunning       = "RUNNING"
	StatusCompleted     = "COMPLETED"
	StatusFailed        = "FAILED"
	StatusPartialFailed = "PARTIAL_FAILED"
	StatusCancelled     = "CANCELLED"
)

const (
	SyncTypeAdAccounts       = "ad_accounts"
	SyncTypeCampaigns        = "campaigns"
	SyncTypeAdsets           = "adsets"
	SyncTypeAds              = "ads"
	SyncTypeAdCreatives      = "ad_creatives"
	SyncTypeCampaignInsights = "campaign_insights"
	SyncTypeAdInsights       = "ad_insights"
	SyncTypeBusinesses       = "businesses"
)

type MetaSyncBatch struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	BatchCode string `gorm:"size:100;not null;uniqueIndex" json:"batch_code"`
	Provider  string `gorm:"size:50;not null;default:meta" json:"provider"`

	AdAccountID   string  `gorm:"size:100;not null;index" json:"ad_account_id"`
	AdAccountName *string `gorm:"size:255" json:"ad_account_name"`

	SyncMode  string `gorm:"size:50;not null;default:scheduled" json:"sync_mode"`
	SyncScope string `gorm:"size:50;not null;default:incremental" json:"sync_scope"`

	DatePreset *string    `gorm:"size:50" json:"date_preset"`
	DateStart  *time.Time `json:"date_start"`
	DateStop   *time.Time `json:"date_stop"`

	Status string `gorm:"size:50;not null;default:PENDING;index" json:"status"`

	StartedAt  *time.Time `gorm:"index" json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	DurationMs uint64     `gorm:"not null;default:0" json:"duration_ms"`

	TotalRecords  uint `gorm:"not null;default:0" json:"total_records"`
	InsertedCount uint `gorm:"not null;default:0" json:"inserted_count"`
	UpdatedCount  uint `gorm:"not null;default:0" json:"updated_count"`
	SkippedCount  uint `gorm:"not null;default:0" json:"skipped_count"`
	FailedCount   uint `gorm:"not null;default:0" json:"failed_count"`

	ProgressPercent uint8 `gorm:"not null;default:0" json:"progress_percent"`

	RequestCount uint `gorm:"not null;default:0" json:"request_count"`
	RateLimitHit bool `gorm:"not null;default:false" json:"rate_limit_hit"`

	ErrorMessage *string        `gorm:"type:text" json:"error_message"`
	Metadata     datatypes.JSON `gorm:"type:json" json:"metadata" swaggertype:"object"`

	Steps []MetaSyncStep `gorm:"foreignKey:BatchID" json:"steps,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (MetaSyncBatch) TableName() string {
	return "meta_sync_batches"
}

type MetaSyncStep struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	BatchID uint64 `gorm:"not null;index" json:"batch_id"`

	SyncType string  `gorm:"size:100;not null;index" json:"sync_type"`
	Endpoint *string `gorm:"size:255" json:"endpoint"`

	Status string `gorm:"size:50;not null;default:PENDING;index" json:"status"`

	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	DurationMs uint64     `gorm:"not null;default:0" json:"duration_ms"`

	TotalRecords  uint `gorm:"not null;default:0" json:"total_records"`
	InsertedCount uint `gorm:"not null;default:0" json:"inserted_count"`
	UpdatedCount  uint `gorm:"not null;default:0" json:"updated_count"`
	SkippedCount  uint `gorm:"not null;default:0" json:"skipped_count"`
	FailedCount   uint `gorm:"not null;default:0" json:"failed_count"`

	RequestCount uint `gorm:"not null;default:0" json:"request_count"`

	CursorBefore *string `gorm:"type:text" json:"cursor_before"`
	CursorAfter  *string `gorm:"type:text" json:"cursor_after"`
	HasNext      bool    `gorm:"not null;default:false" json:"has_next"`

	ErrorCode    *string        `gorm:"size:100" json:"error_code"`
	ErrorMessage *string        `gorm:"type:text" json:"error_message"`
	Metadata     datatypes.JSON `gorm:"type:json" json:"metadata" swaggertype:"object"`

	Batch MetaSyncBatch `gorm:"foreignKey:BatchID" json:"batch,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (MetaSyncStep) TableName() string {
	return "meta_sync_steps"
}

type StepCounts struct {
	TotalRecords  uint
	InsertedCount uint
	UpdatedCount  uint
	SkippedCount  uint
	FailedCount   uint
	RequestCount  uint
}
