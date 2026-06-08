package user

import (
	"github.com/alex/ads_backend/internal/core/user/dto"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
	FindByID(id uint) (*User, error)
	FindAll(filter dto.UserFilter) ([]User, int64, error)
	FindByEmail(email string) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *repository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}

func (r *repository) FindByID(id uint) (*User, error) {
	var user User
	err := r.db.Preload("Roles.Permissions").First(&user, id).Error
	return &user, err
}

func (r *repository) FindAll(filter dto.UserFilter) ([]User, int64, error) {
	var users []User
	var total int64

	q := r.db.Model(&User{}).Scopes(FilterUser(filter))
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

	err := q.Preload("Roles").Limit(limit).Offset((page - 1) * limit).Find(&users).Error
	return users, total, err
}

func (r *repository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Preload("Roles.Permissions").Where("email = ?", email).First(&user).Error
	return &user, err
}
