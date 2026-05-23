package ad_creative

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	FindByCreativeID(creativeID string) (*AdCreative, error)
	Upsert(creative *AdCreative) error
	CreateVersion(version *AdCreativeVersion) error
	FindLatestVersionByCreativeID(creativeID string) (*AdCreativeVersion, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindByCreativeID(creativeID string) (*AdCreative, error) {
	var creative AdCreative
	err := r.db.Where("creative_id = ?", creativeID).First(&creative).Error
	return &creative, err
}

func (r *repository) Upsert(creative *AdCreative) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "creative_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name", "title", "body", "image_url", "video_url",
			"destination_url", "normalized_domain", "url_hash",
			"raw_payload", "synced_at",
		}),
	}).Create(creative).Error
}

func (r *repository) CreateVersion(version *AdCreativeVersion) error {
	return r.db.Create(version).Error
}

func (r *repository) FindLatestVersionByCreativeID(creativeID string) (*AdCreativeVersion, error) {
	var version AdCreativeVersion
	err := r.db.Where("creative_id = ?", creativeID).Order("created_at DESC").First(&version).Error
	return &version, err
}
