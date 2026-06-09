package ad_account

import (
	"github.com/alex/ads_backend/internal/meta/ad_account/dto"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AdAccountFilter struct {
	Search     string
	Page       int
	Limit      int
	BrandID    *uint64
	BusinessID *string
}

type Repository interface {
	UpsertBatch(accounts []MetaAdAccount) error
	FindAll(filter AdAccountFilter) ([]MetaAdAccount, int64, error)
	FindUnassigned(filter AdAccountFilter) ([]MetaAdAccount, int64, error)
	FindByID(id string) (*MetaAdAccount, error)
	GetUniqueBusinesses() ([]dto.BusinessOptionResponse, error)
	Update(account *MetaAdAccount) error
	UpdateBrandID(id string, brandID *uint64) error
	UpdateBrandIDBatch(ids []string, brandID *uint64) error
	UpdateBrandIDByBusiness(businessID string, brandID *uint64) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) UpsertBatch(accounts []MetaAdAccount) error {
	if len(accounts) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"account_status",
			"currency",
			"timezone_name",
			"business_id",
			"business_name",
			"is_active",
			"synced_at",
		}),
	}).CreateInBatches(accounts, 100).Error
}

func (r *repository) FindAll(filter AdAccountFilter) ([]MetaAdAccount, int64, error) {
	var accounts []MetaAdAccount
	var total int64

	query := r.db.Model(&MetaAdAccount{})

	if filter.Search != "" {
		query = query.Where("name LIKE ?", "%"+filter.Search+"%")
	}
	
	if filter.BusinessID != nil {
		query = query.Where("business_id = ?", *filter.BusinessID)
	}

	if filter.BrandID != nil {
		query = query.Where("brand_id = ?", *filter.BrandID)
	}

	query.Count(&total)

	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit

	err := query.Order("name ASC").Limit(filter.Limit).Offset(offset).Find(&accounts).Error
	return accounts, total, err
}

func (r *repository) FindByID(id string) (*MetaAdAccount, error) {
	var account MetaAdAccount
	err := r.db.First(&account, "id = ?", id).Error
	return &account, err
}

func (r *repository) Update(account *MetaAdAccount) error {
	return r.db.Save(account).Error
}

func (r *repository) UpdateBrandID(id string, brandID *uint64) error {
	return r.db.Model(&MetaAdAccount{}).Where("id = ?", id).Update("brand_id", brandID).Error
}

func (r *repository) UpdateBrandIDBatch(ids []string, brandID *uint64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Model(&MetaAdAccount{}).Where("id IN ?", ids).Update("brand_id", brandID).Error
}

func (r *repository) GetUniqueBusinesses() ([]dto.BusinessOptionResponse, error) {
	var businesses []dto.BusinessOptionResponse
	err := r.db.Model(&MetaAdAccount{}).
		Select("business_id, business_name").
		Where("business_id IS NOT NULL AND business_id != ''").
		Group("business_id, business_name").
		Find(&businesses).Error
	return businesses, err
}

func (r *repository) UpdateBrandIDByBusiness(businessID string, brandID *uint64) error {
	return r.db.Model(&MetaAdAccount{}).Where("business_id = ?", businessID).Update("brand_id", brandID).Error
}

func (r *repository) FindUnassigned(filter AdAccountFilter) ([]MetaAdAccount, int64, error) {
	var accounts []MetaAdAccount
	var total int64

	query := r.db.Model(&MetaAdAccount{}).Where("brand_id IS NULL")

	if filter.Search != "" {
		query = query.Where("name LIKE ?", "%"+filter.Search+"%")
	}

	if filter.BusinessID != nil {
		query = query.Where("business_id = ?", *filter.BusinessID)
	}

	query.Count(&total)

	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit

	err := query.Order("name ASC").Limit(filter.Limit).Offset(offset).Find(&accounts).Error
	return accounts, total, err
}
