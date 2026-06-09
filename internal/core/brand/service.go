package brand

import (
	"errors"

	"github.com/alex/ads_backend/internal/core/brand/dto"
	"github.com/alex/ads_backend/pkg/utils"
	"github.com/gosimple/slug"
)

type Service interface {
	Create(req dto.CreateBrandRequest) (dto.BrandResponse, error)
	Update(slug string, req dto.UpdateBrandRequest) (dto.BrandResponse, error)
	Delete(slug string) error
	FindBySlug(slug string) (dto.BrandResponse, error)
	FindAll(filter dto.BrandFilter) ([]dto.BrandResponse, int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(req dto.CreateBrandRequest) (dto.BrandResponse, error) {
	// Generate slug from name
	brandSlug := slug.Make(req.Name)

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	brand := &Brand{
		Name:        req.Name,
		Slug:        brandSlug,
		Description: req.Description,
		IsActive:    isActive,
	}

	// Process photo if present
	if req.Photo != nil {
		photoUrl, err := utils.ProcessBrandPhoto(req.Photo, brandSlug)
		if err != nil {
			return dto.BrandResponse{}, err
		}
		brand.Photo = &photoUrl
	}

	if err := s.repo.Create(brand); err != nil {
		return dto.BrandResponse{}, err
	}

	return toBrandResponse(*brand), nil
}

func (s *service) Update(brandSlug string, req dto.UpdateBrandRequest) (dto.BrandResponse, error) {
	brand, err := s.repo.FindBySlug(brandSlug)
	if err != nil {
		return dto.BrandResponse{}, errors.New("brand not found")
	}

	if req.Name != nil {
		brand.Name = *req.Name
		brand.Slug = slug.Make(*req.Name) // Update slug if name changes
	}
	
	if req.Photo != nil {
		photoUrl, err := utils.ProcessBrandPhoto(req.Photo, brand.Slug)
		if err != nil {
			return dto.BrandResponse{}, err
		}
		brand.Photo = &photoUrl
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

func (s *service) Delete(brandSlug string) error {
	_, err := s.repo.FindBySlug(brandSlug)
	if err != nil {
		return errors.New("brand not found")
	}

	return s.repo.DeleteBySlug(brandSlug)
}

func (s *service) FindBySlug(brandSlug string) (dto.BrandResponse, error) {
	brand, err := s.repo.FindBySlug(brandSlug)
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
		ID:             b.ID,
		Slug:           b.Slug,
		Name:           b.Name,
		Photo:          b.Photo,
		Description:    b.Description,
		IsActive:       b.IsActive,
		AdAccountCount: b.AdAccountCount,
		CreatedAt:      b.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      b.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
