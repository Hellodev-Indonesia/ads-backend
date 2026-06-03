package user_test

import (
	"errors"
	"testing"
	"time"

	"github.com/alex/ads_backend/internal/core/role"
	"github.com/alex/ads_backend/internal/core/user"
	"github.com/alex/ads_backend/internal/core/user/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupService(t *testing.T) (*user.MockRepository, *role.MockRepository, user.Service) {
	mockRepo := user.NewMockRepository(t)
	mockRoleRepo := role.NewMockRepository(t)
	svc := user.NewService(mockRepo, mockRoleRepo)
	return mockRepo, mockRoleRepo, svc
}

func TestService_Create(t *testing.T) {
	mockRepo, mockRoleRepo, svc := setupService(t)

	req := dto.UserRequest{
		Name:     "John Doe",
		Email:    "john@test.com",
		Password: "password123",
		RoleIDs:  []uint{1},
	}

	roles := []role.Role{{ID: 1, Name: "Admin"}}
	mockRoleRepo.On("FindByIDs", req.RoleIDs).Return(roles, nil)

	mockRepo.On("Create", mock.AnythingOfType("*user.User")).Run(func(args mock.Arguments) {
		u := args.Get(0).(*user.User)
		u.ID = 1
	}).Return(nil)

	created, err := svc.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, uint(1), created.ID)
	assert.Equal(t, "John Doe", created.Name)
	assert.Len(t, created.Roles, 1)
}

func TestService_Update(t *testing.T) {
	mockRepo, mockRoleRepo, svc := setupService(t)

	req := dto.UserRequest{
		Name:    "John Updated",
		Email:   "john_updated@test.com",
		RoleIDs: []uint{2},
	}

	existingUser := &user.User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@test.com",
	}

	roles := []role.Role{{ID: 2, Name: "User"}}

	mockRepo.On("FindByID", uint(1)).Return(existingUser, nil)
	mockRoleRepo.On("FindByIDs", req.RoleIDs).Return(roles, nil)
	mockRepo.On("Update", mock.AnythingOfType("*user.User")).Return(nil)

	updated, err := svc.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "John Updated", updated.Name)
	assert.Len(t, updated.Roles, 1)
}

func TestService_Update_NotFound(t *testing.T) {
	mockRepo, _, svc := setupService(t)

	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	updated, err := svc.Update(99, dto.UserRequest{})

	assert.Error(t, err)
	assert.Nil(t, updated)
}

func TestService_Delete(t *testing.T) {
	mockRepo, _, svc := setupService(t)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1)

	assert.NoError(t, err)
}

func TestService_FindByID(t *testing.T) {
	mockRepo, _, svc := setupService(t)

	now := time.Now()
	existingUser := &user.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@test.com",
		CreatedAt: now,
		Roles: []role.Role{
			{ID: 1, Name: "Admin"},
		},
	}

	mockRepo.On("FindByID", uint(1)).Return(existingUser, nil)

	resp, err := svc.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, "John Doe", resp.Name)
	assert.Len(t, resp.Roles, 1)
}

func TestService_FindAll(t *testing.T) {
	mockRepo, _, svc := setupService(t)

	filter := dto.UserFilter{
		Page:  1,
		Limit: 10,
	}

	users := []user.User{
		{ID: 1, Name: "John"},
		{ID: 2, Name: "Jane"},
	}

	mockRepo.On("FindAll", filter).Return(users, int64(2), nil)

	resp, meta, err := svc.FindAll(filter)

	assert.NoError(t, err)
	assert.NotNil(t, meta)
	assert.Len(t, resp, 2)
	assert.Equal(t, 2, meta.Total)
}
