package alert

import (
	"github.com/alex/ads_backend/internal/notification/alert/dto"
	"gorm.io/gorm"
)

type Repository interface {
	Create(a *Alert) error
	Update(a *Alert) error
	FindByID(id uint64) (*Alert, error)
	FindAll(filter dto.AlertFilter) ([]Alert, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(a *Alert) error {
	return r.db.Create(a).Error
}

func (r *repository) Update(a *Alert) error {
	return r.db.Save(a).Error
}

func (r *repository) FindByID(id uint64) (*Alert, error) {
	var a Alert
	err := r.db.First(&a, id).Error
	return &a, err
}

func (r *repository) FindAll(filter dto.AlertFilter) ([]Alert, int64, error) {
	var alerts []Alert
	var total int64

	q := r.db.Model(&Alert{})
	if filter.BrandID != nil {
		q = q.Where("brand_id = ?", *filter.BrandID)
	}
	if filter.Severity != "" {
		q = q.Where("severity = ?", filter.Severity)
	}
	if filter.IsRead != nil {
		q = q.Where("is_read = ?", *filter.IsRead)
	}
	if filter.DateStart != "" {
		q = q.Where("created_at >= ?", filter.DateStart)
	}
	if filter.DateStop != "" {
		q = q.Where("created_at <= ?", filter.DateStop+" 23:59:59")
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := filter.Limit
	page := filter.Page
	if limit <= 0 {
		limit = 25
	}
	if page <= 0 {
		page = 1
	}

	err := q.Order("created_at DESC").Limit(limit).Offset((page - 1) * limit).Find(&alerts).Error
	return alerts, total, err
}
