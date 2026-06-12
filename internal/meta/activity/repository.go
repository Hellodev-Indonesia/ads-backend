package activity

import (
	"github.com/alex/ads_backend/internal/meta/activity/dto"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	UpsertBatch(activities []MetaActivity) error
	FindAll(filter dto.ActivityFilter) ([]ActivityWithAdAccount, int64, error)
	FindAllByBrand(brandID uint64, filter dto.ActivityFilter) ([]ActivityWithAdAccount, int64, error)
}

type ActivityWithAdAccount struct {
	MetaActivity
	AdAccountName         *string `gorm:"column:ad_account_name"`
	AdAccountBusinessName *string `gorm:"column:ad_account_business_name"`
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) UpsertBatch(activities []MetaActivity) error {
	if len(activities) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(activities, 100).Error
}

func (r *repository) FindAllByBrand(brandID uint64, filter dto.ActivityFilter) ([]ActivityWithAdAccount, int64, error) {
	var list []ActivityWithAdAccount
	var total int64

	q := r.db.Table("meta_activities act").
		Joins("JOIN meta_ad_accounts a ON act.ad_account_id = a.id").
		Where("a.brand_id = ?", brandID)

	q.Count(&total)

	limit := filter.Limit
	page := filter.Page
	if limit <= 0 {
		limit = 25
	}
	if page <= 0 {
		page = 1
	}

	err := q.Select("act.*, a.name as ad_account_name, a.business_name as ad_account_business_name").
		Order("act.event_time DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Scan(&list).Error

	return list, total, err
}

func (r *repository) FindAll(filter dto.ActivityFilter) ([]ActivityWithAdAccount, int64, error) {
	var list []ActivityWithAdAccount
	var total int64

	q := r.db.Table("meta_activities act").
		Joins("JOIN meta_ad_accounts a ON act.ad_account_id = a.id")

	q.Count(&total)

	limit := filter.Limit
	page := filter.Page
	if limit <= 0 {
		limit = 25
	}
	if page <= 0 {
		page = 1
	}

	err := q.Select("act.*, a.name as ad_account_name, a.business_name as ad_account_business_name").
		Order("act.event_time DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Scan(&list).Error

	return list, total, err
}
