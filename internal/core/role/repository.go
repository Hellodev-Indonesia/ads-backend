package role

import (
	"github.com/alex/ads_backend/internal/core/permission"
	"gorm.io/gorm"
)

type Repository interface {
	Create(role *Role) error
	Update(role *Role) error
	Delete(id uint) error
	FindAll() ([]Role, error)
	FindByID(id uint) (*Role, error)
	AssignPermissions(role *Role, permissions []permission.Permission) error
	FindByIDs(ids []uint) ([]Role, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(role *Role) error {
	return r.db.Create(role).Error
}

func (r *repository) Update(role *Role) error {
	return r.db.Save(role).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Role{}, id).Error
}

func (r *repository) FindAll() ([]Role, error) {
	var roles []Role
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *repository) FindByID(id uint) (*Role, error) {
	var role Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	return &role, err
}

func (r *repository) AssignPermissions(role *Role, permissions []permission.Permission) error {
	return r.db.Model(role).Association("Permissions").Replace(permissions)
}

func (r *repository) FindByIDs(ids []uint) ([]Role, error) {
	var roles []Role
	err := r.db.Where("id IN ?", ids).Find(&roles).Error
	return roles, err
}
