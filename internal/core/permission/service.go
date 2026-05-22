package permission

import (
	"github.com/alex/ads_backend/internal/core/permission/dto"
	"github.com/alex/ads_backend/pkg/response"
)

type Service interface {
	Create(req dto.PermissionRequest) (*Permission, error)
	Update(id uint, req dto.PermissionRequest) (*Permission, error)
	Delete(id uint) error
	FindByID(id uint) (*dto.PermissionResponse, error)
	FindAll(filter dto.PermissionFilter) ([]dto.PermissionResponse, *response.PaginationMeta, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(req dto.PermissionRequest) (*Permission, error) {
	permission := &Permission{
		Name:        req.Name,
		Description: req.Description,
	}
	err := s.repo.Create(permission)
	return permission, err
}

func (s *service) Update(id uint, req dto.PermissionRequest) (*Permission, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	p.Name = req.Name
	p.Description = req.Description

	err = s.repo.Update(p)
	return p, err
}

func (s *service) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *service) FindByID(id uint) (*dto.PermissionResponse, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.PermissionResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
	}, nil
}

func (s *service) FindAll(filter dto.PermissionFilter) ([]dto.PermissionResponse, *response.PaginationMeta, error) {
	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	permissions, total, err := s.repo.FindPaginated(filter)
	if err != nil {
		return nil, nil, err
	}

	var resp []dto.PermissionResponse
	for _, p := range permissions {
		resp = append(resp, dto.PermissionResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
		})
	}

	lastPage := int(total) / filter.Limit
	if int(total)%filter.Limit > 0 {
		lastPage++
	}

	meta := &response.PaginationMeta{
		Page:     filter.Page,
		Limit:    filter.Limit,
		Total:    int(total),
		LastPage: lastPage,
	}
	return resp, meta, nil
}
