package permission_test

import (
	"errors"
	"testing"

	"github.com/alex/ads_backend/internal/core/permission"
	"github.com/alex/ads_backend/internal/core/permission/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupService(t *testing.T) (*permission.MockRepository, permission.Service) {
	mockRepo := permission.NewMockRepository(t)
	svc := permission.NewService(mockRepo)
	return mockRepo, svc
}

func TestService_Create(t *testing.T) {
	mockRepo, svc := setupService(t)

	req := dto.PermissionRequest{
		Name:        "core.user.view",
		Description: "View user permission",
	}

	mockRepo.On("Create", mock.AnythingOfType("*permission.Permission")).Run(func(args mock.Arguments) {
		p := args.Get(0).(*permission.Permission)
		p.ID = 1
	}).Return(nil)

	created, err := svc.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, uint(1), created.ID)
	assert.Equal(t, "core.user.view", created.Name)
}

func TestService_Update(t *testing.T) {
	mockRepo, svc := setupService(t)

	req := dto.PermissionRequest{
		Name:        "core.user.edit",
		Description: "Edit user permission",
	}

	existing := &permission.Permission{
		ID:          1,
		Name:        "core.user.view",
		Description: "Old desc",
	}

	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*permission.Permission")).Return(nil)

	updated, err := svc.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "core.user.edit", updated.Name)
	assert.Equal(t, "Edit user permission", updated.Description)
}

func TestService_Update_NotFound(t *testing.T) {
	mockRepo, svc := setupService(t)

	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	updated, err := svc.Update(99, dto.PermissionRequest{})

	assert.Error(t, err)
	assert.Nil(t, updated)
}

func TestService_Delete(t *testing.T) {
	mockRepo, svc := setupService(t)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1)

	assert.NoError(t, err)
}

func TestService_FindByID(t *testing.T) {
	mockRepo, svc := setupService(t)

	existing := &permission.Permission{
		ID:          1,
		Name:        "core.user.view",
		Description: "View user",
	}

	mockRepo.On("FindByID", uint(1)).Return(existing, nil)

	resp, err := svc.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, "core.user.view", resp.Name)
}

func TestService_FindAll(t *testing.T) {
	mockRepo, svc := setupService(t)

	filter := dto.PermissionFilter{
		Page:  1,
		Limit: 10,
	}

	perms := []permission.Permission{
		{ID: 1, Name: "core.user.view"},
		{ID: 2, Name: "core.user.edit"},
	}

	mockRepo.On("FindPaginated", filter).Return(perms, int64(2), nil)

	resp, meta, err := svc.FindAll(filter)

	assert.NoError(t, err)
	assert.NotNil(t, meta)
	assert.Len(t, resp, 2)
	assert.Equal(t, 2, meta.Total)
}
