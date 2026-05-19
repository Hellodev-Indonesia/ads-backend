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
	FindAll(filter dto.UserFilter) ([]User, error)
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

func (r *repository) FindAll(filter dto.UserFilter) ([]User, error) {
	var users []User
	err := r.db.Scopes(FilterUser(filter)).Preload("Roles").Find(&users).Error
	return users, err
}

func (r *repository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Preload("Roles.Permissions").Where("email = ?", email).First(&user).Error
	return &user, err
}
