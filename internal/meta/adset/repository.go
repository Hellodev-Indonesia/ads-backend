package adset

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AdSetFilter struct {
	AccountID  string
	BrandID    *uint64
	CampaignID string
	Status     string
	Search     string
	Page       int
	Limit      int
}

type Repository interface {
	Upsert(adset *MetaAdSet) error
	UpsertBatch(adsets []MetaAdSet) error
	FindAll(filter AdSetFilter) ([]MetaAdSet, int64, error)
	FindByCampaignID(campaignID string) ([]MetaAdSet, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Upsert(adset *MetaAdSet) error {
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(adset).Error
}

func (r *repository) UpsertBatch(adsets []MetaAdSet) error {
	if len(adsets) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(adsets, 100).Error
}

func (r *repository) FindAll(filter AdSetFilter) ([]MetaAdSet, int64, error) {
	var adsets []MetaAdSet
	var total int64

	query := r.db.Model(&MetaAdSet{})

	if filter.CampaignID != "" {
		query = query.Where("meta_ad_sets.campaign_id = ?", filter.CampaignID)
	}
	if filter.Status != "" {
		query = query.Where("meta_ad_sets.status = ?", filter.Status)
	}
	if filter.Search != "" {
		query = query.Where("meta_ad_sets.name LIKE ?", "%"+filter.Search+"%")
	}
	if filter.BrandID != nil {
		query = query.Joins("JOIN meta_campaigns ON meta_ad_sets.campaign_id = meta_campaigns.id").
			Joins("JOIN meta_ad_accounts ON meta_campaigns.account_id = meta_ad_accounts.id").
			Where("meta_ad_accounts.brand_id = ?", *filter.BrandID)
	}

	query.Count(&total)

	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit

	err := query.Order("created_time DESC").Limit(filter.Limit).Offset(offset).Find(&adsets).Error
	return adsets, total, err
}

func (r *repository) FindByCampaignID(campaignID string) ([]MetaAdSet, error) {
	var adsets []MetaAdSet
	err := r.db.Where("campaign_id = ?", campaignID).Find(&adsets).Error
	return adsets, err
}
