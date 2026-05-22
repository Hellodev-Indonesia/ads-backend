package ad_account

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AdAccountFilter struct {
	Search string
	Page   int
	Limit  int
}

type Repository interface {
	UpsertBatch(accounts []MetaAdAccount) error
	FindAll(filter AdAccountFilter) ([]MetaAdAccount, int64, error)
	FindUnassigned(filter AdAccountFilter) ([]MetaAdAccount, int64, error)
	FindByID(id string) (*MetaAdAccount, error)
	Update(account *MetaAdAccount) error
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
		UpdateAll: true,
	}).CreateInBatches(accounts, 100).Error
}

func (r *repository) FindAll(filter AdAccountFilter) ([]MetaAdAccount, int64, error) {
	var accounts []MetaAdAccount
	var total int64

	query := r.db.Model(&MetaAdAccount{})

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

func (r *repository) FindUnassigned(filter AdAccountFilter) ([]MetaAdAccount, int64, error) {
	var accounts []MetaAdAccount
	var total int64

	query := r.db.Model(&MetaAdAccount{}).Where("brand_id IS NULL")

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

	err := query.Order("name ASC").Limit(filter.Limit).Offset(offset).Find(&accounts).Error
	return accounts, total, err
}
