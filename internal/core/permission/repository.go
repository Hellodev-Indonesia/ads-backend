package permission

import (
	"github.com/alex/ads_backend/internal/core/permission/dto"
	"gorm.io/gorm"
)

type Repository interface {
	Create(permission *Permission) error
	Update(permission *Permission) error
	Delete(id uint) error
	FindByID(id uint) (*Permission, error)
	FindAll() ([]Permission, error)
	FindPaginated(filter dto.PermissionFilter) ([]Permission, int64, error)
	FindByIDs(ids []uint) ([]Permission, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(permission *Permission) error {
	return r.db.Create(permission).Error
}

func (r *repository) Update(permission *Permission) error {
	return r.db.Save(permission).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Permission{}, id).Error
}

func (r *repository) FindByID(id uint) (*Permission, error) {
	var p Permission
	err := r.db.First(&p, id).Error
	return &p, err
}

func (r *repository) FindAll() ([]Permission, error) {
	var permissions []Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}

func (r *repository) FindPaginated(filter dto.PermissionFilter) ([]Permission, int64, error) {
	var permissions []Permission
	var total int64

	q := r.db.Model(&Permission{})
	if filter.Name != "" {
		q = q.Where("name LIKE ?", "%"+filter.Name+"%")
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

	err := q.Limit(limit).Offset((page - 1) * limit).Find(&permissions).Error
	return permissions, total, err
}

func (r *repository) FindByIDs(ids []uint) ([]Permission, error) {
	var permissions []Permission
	err := r.db.Where("id IN ?", ids).Find(&permissions).Error
	return permissions, err
}
