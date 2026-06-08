package brand

import (
	"github.com/alex/ads_backend/internal/core/brand/dto"
	"gorm.io/gorm"
)

type Repository interface {
	Create(brand *Brand) error
	Update(brand *Brand) error
	Delete(id uint64) error
	FindByID(id uint64) (*Brand, error)
	FindAll(filter dto.BrandFilter) ([]Brand, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func FilterBrand(filter dto.BrandFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter.Name != "" {
			db = db.Where("name LIKE ?", "%"+filter.Name+"%")
		}
		if filter.IsActive != nil {
			db = db.Where("is_active = ?", *filter.IsActive)
		}
		return db
	}
}

func (r *repository) Create(brand *Brand) error {
	return r.db.Create(brand).Error
}

func (r *repository) Update(brand *Brand) error {
	return r.db.Save(brand).Error
}

func (r *repository) Delete(id uint64) error {
	return r.db.Delete(&Brand{}, id).Error
}

func (r *repository) FindByID(id uint64) (*Brand, error) {
	var brand Brand
	err := r.db.First(&brand, id).Error
	return &brand, err
}

func (r *repository) FindAll(filter dto.BrandFilter) ([]Brand, int64, error) {
	var brands []Brand
	var total int64

	q := r.db.Model(&Brand{}).Scopes(FilterBrand(filter))
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

	err := q.Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&brands).Error
	return brands, total, err
}
