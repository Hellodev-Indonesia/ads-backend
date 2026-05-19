package role

import (
	"github.com/alex/ads_backend/internal/core/permission"
	permDto "github.com/alex/ads_backend/internal/core/permission/dto"
	"github.com/alex/ads_backend/internal/core/role/dto"
)

type Service interface {
	Create(req dto.RoleRequest) (*Role, error)
	Update(id uint, req dto.RoleRequest) (*Role, error)
	Delete(id uint) error
	FindAll() ([]dto.RoleResponse, error)
	FindByID(id uint) (*dto.RoleResponse, error)
	AssignPermissions(roleID uint, req dto.AssignPermissionRequest) error
}

type service struct {
	repo     Repository
	permRepo permission.Repository
}

func NewService(repo Repository, permRepo permission.Repository) Service {
	return &service{repo, permRepo}
}

func (s *service) Create(req dto.RoleRequest) (*Role, error) {
	role := &Role{
		Name:        req.Name,
		Description: req.Description,
	}
	err := s.repo.Create(role)
	return role, err
}

func (s *service) Update(id uint, req dto.RoleRequest) (*Role, error) {
	r, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	r.Name = req.Name
	r.Description = req.Description

	err = s.repo.Update(r)
	return r, err
}

func (s *service) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *service) FindAll() ([]dto.RoleResponse, error) {
	roles, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var resp []dto.RoleResponse
	for _, r := range roles {
		var perms []permDto.PermissionResponse
		for _, p := range r.Permissions {
			perms = append(perms, permDto.PermissionResponse{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
			})
		}
		resp = append(resp, dto.RoleResponse{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Permissions: perms,
		})
	}
	return resp, nil
}

func (s *service) FindByID(id uint) (*dto.RoleResponse, error) {
	r, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	var perms []permDto.PermissionResponse
	for _, p := range r.Permissions {
		perms = append(perms, permDto.PermissionResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
		})
	}

	return &dto.RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: perms,
	}, nil
}

func (s *service) AssignPermissions(roleID uint, req dto.AssignPermissionRequest) error {
	role, err := s.repo.FindByID(roleID)
	if err != nil {
		return err
	}

	permissions, err := s.permRepo.FindByIDs(req.PermissionIDs)
	if err != nil {
		return err
	}

	return s.repo.AssignPermissions(role, permissions)
}
