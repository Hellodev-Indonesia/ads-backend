package role_test

import (
	"errors"
	"testing"

	"github.com/alex/ads_backend/internal/core/permission"
	"github.com/alex/ads_backend/internal/core/role"
	"github.com/alex/ads_backend/internal/core/role/dto"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupService is a helper to set up the service with mocks
func setupService(t *testing.T) (*role.MockRepository, *permission.MockRepository, role.Service) {
	mockRoleRepo := role.NewMockRepository(t)
	mockPermRepo := permission.NewMockRepository(t)
	svc := role.NewService(mockRoleRepo, mockPermRepo)
	return mockRoleRepo, mockPermRepo, svc
}

func TestService_Create(t *testing.T) {
	mockRoleRepo, _, svc := setupService(t)

	req := dto.RoleRequest{
		Name:        "Admin",
		Description: "Administrator role",
	}

	mockRoleRepo.On("Create", mock.AnythingOfType("*role.Role")).Run(func(args mock.Arguments) {
		r := args.Get(0).(*role.Role)
		r.ID = 1
	}).Return(nil)

	createdRole, err := svc.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, createdRole)
	assert.Equal(t, "Admin", createdRole.Name)
	assert.Equal(t, uint(1), createdRole.ID)
}

func TestService_Update(t *testing.T) {
	mockRoleRepo, _, svc := setupService(t)

	req := dto.RoleRequest{
		Name:        "Super Admin",
		Description: "Updated role",
	}

	existingRole := &role.Role{
		ID:          1,
		Name:        "Admin",
		Description: "Old role",
	}

	mockRoleRepo.On("FindByID", uint(1)).Return(existingRole, nil)
	mockRoleRepo.On("Update", mock.AnythingOfType("*role.Role")).Return(nil)

	updatedRole, err := svc.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, updatedRole)
	assert.Equal(t, "Super Admin", updatedRole.Name)
	assert.Equal(t, "Updated role", updatedRole.Description)
}

func TestService_Update_NotFound(t *testing.T) {
	mockRoleRepo, _, svc := setupService(t)

	mockRoleRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	updatedRole, err := svc.Update(99, dto.RoleRequest{})

	assert.Error(t, err)
	assert.Nil(t, updatedRole)
}

func TestService_Delete(t *testing.T) {
	mockRoleRepo, _, svc := setupService(t)

	mockRoleRepo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1)

	assert.NoError(t, err)
}

func TestService_FindByID(t *testing.T) {
	mockRoleRepo, _, svc := setupService(t)

	existingRole := &role.Role{
		ID:          1,
		Name:        "Admin",
		Description: "Admin role",
		Permissions: []permission.Permission{
			{ID: 1, Name: "core.user.view", Description: "View users"},
		},
	}

	mockRoleRepo.On("FindByID", uint(1)).Return(existingRole, nil)

	resp, err := svc.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint(1), resp.ID)
	assert.Len(t, resp.Permissions, 1)
	assert.Equal(t, "core.user.view", resp.Permissions[0].Name)
}

func TestService_FindAll(t *testing.T) {
	mockRoleRepo, _, svc := setupService(t)

	filter := dto.RoleFilter{
		Page:  1,
		Limit: 10,
	}

	roles := []role.Role{
		{ID: 1, Name: "Admin"},
		{ID: 2, Name: "User"},
	}

	mockRoleRepo.On("FindAll", filter).Return(roles, int64(2), nil)

	resp, meta, err := svc.FindAll(filter)

	assert.NoError(t, err)
	assert.NotNil(t, meta)
	assert.Len(t, resp, 2)
	assert.Equal(t, 2, meta.Total)
	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 10, meta.Limit)
}

func TestService_AssignPermissions(t *testing.T) {
	mockRoleRepo, mockPermRepo, svc := setupService(t)

	existingRole := &role.Role{ID: 1, Name: "Admin"}
	perms := []permission.Permission{
		{ID: 1, Name: "core.user.view"},
		{ID: 2, Name: "core.user.create"},
	}

	req := dto.AssignPermissionRequest{
		PermissionIDs: []uint{1, 2},
	}

	mockRoleRepo.On("FindByID", uint(1)).Return(existingRole, nil)
	mockPermRepo.On("FindByIDs", req.PermissionIDs).Return(perms, nil)
	mockRoleRepo.On("AssignPermissions", existingRole, perms).Return(nil)

	err := svc.AssignPermissions(1, req)

	assert.NoError(t, err)
}
