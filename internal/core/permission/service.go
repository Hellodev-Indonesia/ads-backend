package permission

import "github.com/alex/ads_backend/internal/core/permission/dto"

type Service interface {
	Create(req dto.PermissionRequest) (*Permission, error)
	Update(id uint, req dto.PermissionRequest) (*Permission, error)
	Delete(id uint) error
	FindByID(id uint) (*dto.PermissionResponse, error)
	FindAll() ([]dto.PermissionResponse, error)
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

func (s *service) FindAll() ([]dto.PermissionResponse, error) {
	permissions, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var resp []dto.PermissionResponse
	for _, p := range permissions {
		resp = append(resp, dto.PermissionResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
		})
	}
	return resp, nil
}
