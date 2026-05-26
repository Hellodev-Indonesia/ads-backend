package ads

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AdFilter struct {
	AccountID  string
	CampaignID string
	AdSetID    string
	Status     string
	Search     string
	Page       int
	Limit      int
}

type Repository interface {
	Upsert(ad *MetaAd) error
	UpsertBatch(ads []MetaAd) error
	FindAll(filter AdFilter) ([]MetaAd, int64, error)
	FindCreativeRawPayload(creativeID string) (string, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Upsert(ad *MetaAd) error {
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(ad).Error
}

func (r *repository) UpsertBatch(ads []MetaAd) error {
	if len(ads) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(ads, 100).Error
}

func (r *repository) FindAll(filter AdFilter) ([]MetaAd, int64, error) {
	var adsList []MetaAd
	var total int64

	query := r.db.Model(&MetaAd{})

	if filter.CampaignID != "" {
		query = query.Where("campaign_id = ?", filter.CampaignID)
	}
	if filter.AdSetID != "" {
		query = query.Where("adset_id = ?", filter.AdSetID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		query = query.Where("name LIKE ?", "%"+filter.Search+"%")
	}

	query.Count(&total)

	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit

	err := query.Order("created_time DESC").Limit(filter.Limit).Offset(offset).Find(&adsList).Error
	return adsList, total, err
}

func (r *repository) FindCreativeRawPayload(creativeID string) (string, error) {
	var result struct {
		RawPayload string
	}
	err := r.db.Table("ad_creatives").Select("raw_payload").Where("creative_id = ? AND deleted_at IS NULL", creativeID).First(&result).Error
	return result.RawPayload, err
}
