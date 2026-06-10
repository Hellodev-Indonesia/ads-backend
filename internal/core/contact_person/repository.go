package contact_person

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(limit, offset int) ([]ContactPerson, error)
	Count() (int64, error)
	FindByID(id uint64) (*ContactPerson, error)
	Create(contactPerson *ContactPerson) error
	Update(contactPerson *ContactPerson) error
	Delete(id uint64) error
}

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) FindAll(limit, offset int) ([]ContactPerson, error) {
	var contactPersons []ContactPerson
	err := r.db.Limit(limit).Offset(offset).Order("id DESC").Find(&contactPersons).Error
	return contactPersons, err
}

func (r *repositoryImpl) Count() (int64, error) {
	var count int64
	err := r.db.Model(&ContactPerson{}).Count(&count).Error
	return count, err
}

func (r *repositoryImpl) FindByID(id uint64) (*ContactPerson, error) {
	var contactPerson ContactPerson
	if err := r.db.First(&contactPerson, id).Error; err != nil {
		return nil, err
	}
	return &contactPerson, nil
}

func (r *repositoryImpl) Create(contactPerson *ContactPerson) error {
	return r.db.Create(contactPerson).Error
}

func (r *repositoryImpl) Update(contactPerson *ContactPerson) error {
	return r.db.Save(contactPerson).Error
}

func (r *repositoryImpl) Delete(id uint64) error {
	return r.db.Delete(&ContactPerson{}, id).Error
}
