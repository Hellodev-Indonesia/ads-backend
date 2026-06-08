package brand

import (
	"errors"

	"github.com/alex/ads_backend/internal/core/brand/dto"
)

type Service interface {
	Create(req dto.CreateBrandRequest) (dto.BrandResponse, error)
	Update(id uint64, req dto.UpdateBrandRequest) (dto.BrandResponse, error)
	Delete(id uint64) error
	FindByID(id uint64) (dto.BrandResponse, error)
	FindAll(filter dto.BrandFilter) ([]dto.BrandResponse, int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(req dto.CreateBrandRequest) (dto.BrandResponse, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	brand := &Brand{
		Name:        req.Name,
		Photo:       req.Photo,
		Description: req.Description,
		IsActive:    isActive,
	}

	if err := s.repo.Create(brand); err != nil {
		return dto.BrandResponse{}, err
	}

	return toBrandResponse(*brand), nil
}

func (s *service) Update(id uint64, req dto.UpdateBrandRequest) (dto.BrandResponse, error) {
	brand, err := s.repo.FindByID(id)
	if err != nil {
		return dto.BrandResponse{}, errors.New("brand not found")
	}

	if req.Name != nil {
		brand.Name = *req.Name
	}
	if req.Photo != nil {
		brand.Photo = req.Photo
	}
	if req.Description != nil {
		brand.Description = req.Description
	}
	if req.IsActive != nil {
		brand.IsActive = *req.IsActive
	}

	if err := s.repo.Update(brand); err != nil {
		return dto.BrandResponse{}, err
	}

	return toBrandResponse(*brand), nil
}

func (s *service) Delete(id uint64) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("brand not found")
	}

	return s.repo.Delete(id)
}

func (s *service) FindByID(id uint64) (dto.BrandResponse, error) {
	brand, err := s.repo.FindByID(id)
	if err != nil {
		return dto.BrandResponse{}, errors.New("brand not found")
	}

	return toBrandResponse(*brand), nil
}

func (s *service) FindAll(filter dto.BrandFilter) ([]dto.BrandResponse, int64, error) {
	brands, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.BrandResponse
	for _, b := range brands {
		responses = append(responses, toBrandResponse(b))
	}

	return responses, total, nil
}

func toBrandResponse(b Brand) dto.BrandResponse {
	return dto.BrandResponse{
		ID:          b.ID,
		Name:        b.Name,
		Photo:       b.Photo,
		Description: b.Description,
		IsActive:    b.IsActive,
		CreatedAt:   b.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   b.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
