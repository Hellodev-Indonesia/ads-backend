package campaign

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CampaignFilter struct {
	Status    string
	AccountID string
	Search    string
	Page      int
	Limit     int
}

type Repository interface {
	Upsert(campaign *MetaCampaign) error
	UpsertBatch(campaigns []MetaCampaign) error
	FindAll(filter CampaignFilter) ([]MetaCampaign, int64, error)
	FindByID(id string) (*MetaCampaign, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Upsert(campaign *MetaCampaign) error {
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(campaign).Error
}

func (r *repository) UpsertBatch(campaigns []MetaCampaign) error {
	if len(campaigns) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(campaigns, 100).Error
}

func (r *repository) FindAll(filter CampaignFilter) ([]MetaCampaign, int64, error) {
	var campaigns []MetaCampaign
	var total int64

	query := r.db.Model(&MetaCampaign{})

	if filter.AccountID != "" {
		query = query.Where("account_id = ?", filter.AccountID)
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

	err := query.Order("created_time DESC").Limit(filter.Limit).Offset(offset).Find(&campaigns).Error
	return campaigns, total, err
}

func (r *repository) FindByID(id string) (*MetaCampaign, error) {
	var campaign MetaCampaign
	err := r.db.First(&campaign, "id = ?", id).Error
	return &campaign, err
}
