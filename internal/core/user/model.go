package user

import (
	"time"

	"github.com/alex/ads_backend/internal/core/role"
	"github.com/alex/ads_backend/internal/core/user/dto"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Email     string         `gorm:"size:255;not null;unique" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Roles     []role.Role    `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func FilterUser(filter dto.UserFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter.Name != "" {
			db = db.Where("name LIKE ?", "%"+filter.Name+"%")
		}
		if filter.Email != "" {
			db = db.Where("email LIKE ?", "%"+filter.Email+"%")
		}
		return db
	}
}
