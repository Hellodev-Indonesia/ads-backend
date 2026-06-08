package brand_test

import (
	"errors"
	"testing"
	"time"

	"github.com/alex/ads_backend/internal/core/brand"
	"github.com/alex/ads_backend/internal/core/brand/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupService(t *testing.T) (*brand.MockRepository, brand.Service) {
	mockRepo := brand.NewMockRepository(t)
	svc := brand.NewService(mockRepo)
	return mockRepo, svc
}

func TestService_Create(t *testing.T) {
	mockRepo, svc := setupService(t)

	isActive := true
	req := dto.CreateBrandRequest{
		Name:     "Test Brand",
		IsActive: &isActive,
	}

	mockRepo.On("Create", mock.AnythingOfType("*brand.Brand")).Run(func(args mock.Arguments) {
		b := args.Get(0).(*brand.Brand)
		b.ID = 1
		b.CreatedAt = time.Now()
		b.UpdatedAt = time.Now()
	}).Return(nil)

	resp, err := svc.Create(req)

	assert.NoError(t, err)
	assert.Equal(t, uint64(1), resp.ID)
	assert.Equal(t, "Test Brand", resp.Name)
	assert.True(t, resp.IsActive)
}

func TestService_Update(t *testing.T) {
	mockRepo, svc := setupService(t)

	newName := "Updated Brand"
	req := dto.UpdateBrandRequest{
		Name: &newName,
	}

	existing := &brand.Brand{
		ID:        1,
		Name:      "Old Brand",
		IsActive:  true,
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByID", uint64(1)).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*brand.Brand")).Return(nil)

	resp, err := svc.Update(1, req)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Brand", resp.Name)
}

func TestService_Update_NotFound(t *testing.T) {
	mockRepo, svc := setupService(t)

	mockRepo.On("FindByID", uint64(99)).Return(nil, errors.New("not found"))

	resp, err := svc.Update(99, dto.UpdateBrandRequest{})

	assert.Error(t, err)
	assert.Equal(t, uint64(0), resp.ID)
}

func TestService_Delete(t *testing.T) {
	mockRepo, svc := setupService(t)

	mockRepo.On("FindByID", uint64(1)).Return(&brand.Brand{ID: 1}, nil)
	mockRepo.On("Delete", uint64(1)).Return(nil)

	err := svc.Delete(1)
	assert.NoError(t, err)
}

func TestService_FindByID(t *testing.T) {
	mockRepo, svc := setupService(t)

	existing := &brand.Brand{
		ID:        1,
		Name:      "Brand One",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByID", uint64(1)).Return(existing, nil)

	resp, err := svc.FindByID(1)

	assert.NoError(t, err)
	assert.Equal(t, uint64(1), resp.ID)
	assert.Equal(t, "Brand One", resp.Name)
}

func TestService_FindAll(t *testing.T) {
	mockRepo, svc := setupService(t)

	filter := dto.BrandFilter{
		Page:  1,
		Limit: 10,
	}

	brands := []brand.Brand{
		{ID: 1, Name: "Brand 1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	mockRepo.On("FindAll", filter).Return(brands, int64(1), nil)

	resp, total, err := svc.FindAll(filter)

	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, int64(1), total)
}
