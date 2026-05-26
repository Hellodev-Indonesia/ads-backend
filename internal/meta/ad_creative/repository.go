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
	var creatives []AdCreative
	err := r.db.Where("creative_id = ?", creativeID).Limit(1).Find(&creatives).Error
	if err != nil {
		return nil, err
	}
	if len(creatives) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &creatives[0], nil
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
	var versions []AdCreativeVersion
	err := r.db.Where("creative_id = ?", creativeID).Order("created_at DESC").Limit(1).Find(&versions).Error
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &versions[0], nil
}
