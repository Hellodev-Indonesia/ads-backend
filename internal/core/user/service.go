package user

import (
	"github.com/alex/ads_backend/internal/core/role"
	"github.com/alex/ads_backend/internal/core/user/dto"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(req dto.UserRequest) (*User, error)
	Update(id uint, req dto.UserRequest) (*User, error)
	Delete(id uint) error
	FindAll(filter dto.UserFilter) ([]dto.UserResponse, error)
	FindByID(id uint) (*dto.UserResponse, error)
}

type service struct {
	repo     Repository
	roleRepo role.Repository
}

func NewService(repo Repository, roleRepo role.Repository) Service {
	return &service{repo, roleRepo}
}

func (s *service) Create(req dto.UserRequest) (*User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	
	roles, err := s.roleRepo.FindByIDs(req.RoleIDs)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Roles:    roles,
	}

	err = s.repo.Create(user)
	return user, err
}

func (s *service) Update(id uint, req dto.UserRequest) (*User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	
	if len(req.RoleIDs) > 0 {
		roles, err := s.roleRepo.FindByIDs(req.RoleIDs)
		if err != nil {
			return nil, err
		}
		user.Roles = roles
	}

	err = s.repo.Update(user)
	return user, err
}

func (s *service) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *service) FindAll(filter dto.UserFilter) ([]dto.UserResponse, error) {
	users, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	var resp []dto.UserResponse
	for _, u := range users {
		var roles []dto.RoleBrief
		for _, r := range u.Roles {
			roles = append(roles, dto.RoleBrief{ID: r.ID, Name: r.Name})
		}
		resp = append(resp, dto.UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Roles:     roles,
			CreatedAt: u.CreatedAt,
		})
	}
	return resp, nil
}

func (s *service) FindByID(id uint) (*dto.UserResponse, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	var roles []dto.RoleBrief
	for _, r := range u.Roles {
		roles = append(roles, dto.RoleBrief{ID: r.ID, Name: r.Name})
	}

	return &dto.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Roles:     roles,
		CreatedAt: u.CreatedAt,
	}, nil
}
