package auth

import (
	"errors"

	"github.com/alex/ads_backend/internal/core/auth/dto"
	"github.com/alex/ads_backend/internal/core/permission"
	"github.com/alex/ads_backend/internal/core/user"
	"github.com/alex/ads_backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
}

type service struct {
	userRepo user.Repository
	permRepo permission.Repository
}

func NewService(userRepo user.Repository, permRepo permission.Repository) Service {
	return &service{userRepo, permRepo}
}

func (s *service) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	u, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	var roles []string
	var permissions []string
	isSuperAdmin := false

	for _, role := range u.Roles {
		roles = append(roles, role.Name)
		if role.Name == "Super Admin" {
			isSuperAdmin = true
		}
		if !isSuperAdmin {
			for _, perm := range role.Permissions {
				permissions = append(permissions, perm.Name)
			}
		}
	}

	// If Super Admin, get all permissions from system
	if isSuperAdmin {
		allPerms, err := s.permRepo.FindAll()
		if err == nil {
			permissions = []string{} // Clear existing if any
			for _, p := range allPerms {
				permissions = append(permissions, p.Name)
			}
		}
	}

	token, err := utils.GenerateToken(u.ID, u.Email, roles, permissions)
	if err != nil {
		return nil, err
	}

	centrifugoToken, _ := utils.GenerateCentrifugoToken(string(rune(u.ID)))

	userResp := dto.AuthUserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}

	return &dto.LoginResponse{
		Token:           token,
		CentrifugoToken: centrifugoToken,
		User:            userResp,
		Roles:           roles,
		Permissions:     permissions,
	}, nil
}
