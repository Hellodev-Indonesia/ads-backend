package seeders

import (
	"github.com/alex/ads_backend/internal/core/permission"
	"github.com/alex/ads_backend/internal/core/role"
	"github.com/alex/ads_backend/internal/core/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedCore(db *gorm.DB) {
	// 1. Seed Permissions
	permissions := []permission.Permission{
		{Name: "core.user.view", Description: "View users"},
		{Name: "core.user.create", Description: "Create user"},
		{Name: "core.user.update", Description: "Update user"},
		{Name: "core.user.delete", Description: "Delete user"},
		{Name: "rbac.role.manage", Description: "Manage roles"},
		{Name: "meta.campaign.view", Description: "View campaigns"},
		{Name: "meta.campaign.sync", Description: "Sync campaigns from Meta"},
	}

	for i := range permissions {
		db.FirstOrCreate(&permissions[i], permission.Permission{Name: permissions[i].Name})
	}

	// 2. Seed Roles
	superAdminRole := role.Role{Name: "Super Admin", Description: "Full access to all modules"}
	db.FirstOrCreate(&superAdminRole, role.Role{Name: "Super Admin"})

	adsManagerRole := role.Role{Name: "Ads Manager", Description: "Manage meta ads and insights"}
	db.FirstOrCreate(&adsManagerRole, role.Role{Name: "Ads Manager"})

	// 3. Assign Permissions to Roles (Optional, Super Admin gets all anyway in service)
	// But let's assign them for completeness
	db.Model(&superAdminRole).Association("Permissions").Replace(permissions)

	// 4. Seed Super Admin User
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	adminUser := user.User{
		Name:     "Super Admin",
		Email:    "admin@example.com",
		Password: string(hashedPassword),
	}

	if err := db.Where(user.User{Email: "admin@example.com"}).FirstOrCreate(&adminUser).Error; err == nil {
		// Assign role
		db.Model(&adminUser).Association("Roles").Replace([]role.Role{superAdminRole})
	}
}
