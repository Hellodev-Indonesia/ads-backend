package business

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BusinessFilter struct {
	Search string
	Page   int
	Limit  int
}

type Repository interface {
	UpsertBatch(businesses []MetaBusiness) error
	FindAll(filter BusinessFilter) ([]MetaBusiness, int64, error)
	FindByID(id string) (*MetaBusiness, error)
	GetUniqueBusinessIDsFromAdAccounts() ([]string, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) UpsertBatch(businesses []MetaBusiness) error {
	if len(businesses) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"profile_picture_uri",
			"timezone_id",
			"created_time",
			"synced_at",
		}),
	}).CreateInBatches(businesses, 100).Error
}

func (r *repository) FindAll(filter BusinessFilter) ([]MetaBusiness, int64, error) {
	var businesses []MetaBusiness
	var total int64

	query := r.db.Model(&MetaBusiness{})

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

	err := query.Order("name ASC").Limit(filter.Limit).Offset(offset).Find(&businesses).Error
	return businesses, total, err
}

func (r *repository) FindByID(id string) (*MetaBusiness, error) {
	var business MetaBusiness
	err := r.db.First(&business, "id = ?", id).Error
	return &business, err
}

func (r *repository) GetUniqueBusinessIDsFromAdAccounts() ([]string, error) {
	var ids []string
	err := r.db.Table("meta_ad_accounts").Where("business_id IS NOT NULL AND business_id != ''").Distinct("business_id").Pluck("business_id", &ids).Error
	return ids, err
}
