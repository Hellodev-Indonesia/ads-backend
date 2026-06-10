package contact_person

import (
	"errors"

	"github.com/alex/ads_backend/internal/core/contact_person/dto"
)

var (
	ErrContactPersonNotFound = errors.New("contact person not found")
)

type Service interface {
	FindAll(limit, offset int) ([]dto.ContactPersonListResponse, int64, error)
	FindByID(id uint64) (*dto.ContactPersonResponse, error)
	Create(req dto.ContactPersonRequest) (*dto.ContactPersonResponse, error)
	Update(id uint64, req dto.ContactPersonRequest) (*dto.ContactPersonResponse, error)
	Delete(id uint64) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) FindAll(limit, offset int) ([]dto.ContactPersonListResponse, int64, error) {
	contactPersons, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	var res []dto.ContactPersonListResponse
	for _, cp := range contactPersons {
		res = append(res, dto.ContactPersonListResponse{
			ID:        cp.ID,
			Name:      cp.Name,
			Phone:     cp.Phone,
			CreatedAt: cp.CreatedAt,
		})
	}

	return res, count, nil
}

func (s *serviceImpl) FindByID(id uint64) (*dto.ContactPersonResponse, error) {
	cp, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrContactPersonNotFound
	}

	return &dto.ContactPersonResponse{
		ID:        cp.ID,
		Name:      cp.Name,
		Phone:     cp.Phone,
		CreatedAt: cp.CreatedAt,
		UpdatedAt: cp.UpdatedAt,
	}, nil
}

func (s *serviceImpl) Create(req dto.ContactPersonRequest) (*dto.ContactPersonResponse, error) {
	cp := &ContactPerson{
		Name:  req.Name,
		Phone: req.Phone,
	}

	if err := s.repo.Create(cp); err != nil {
		return nil, err
	}

	return &dto.ContactPersonResponse{
		ID:        cp.ID,
		Name:      cp.Name,
		Phone:     cp.Phone,
		CreatedAt: cp.CreatedAt,
		UpdatedAt: cp.UpdatedAt,
	}, nil
}

func (s *serviceImpl) Update(id uint64, req dto.ContactPersonRequest) (*dto.ContactPersonResponse, error) {
	cp, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrContactPersonNotFound
	}

	cp.Name = req.Name
	cp.Phone = req.Phone

	if err := s.repo.Update(cp); err != nil {
		return nil, err
	}

	return &dto.ContactPersonResponse{
		ID:        cp.ID,
		Name:      cp.Name,
		Phone:     cp.Phone,
		CreatedAt: cp.CreatedAt,
		UpdatedAt: cp.UpdatedAt,
	}, nil
}

func (s *serviceImpl) Delete(id uint64) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return ErrContactPersonNotFound
	}

	return s.repo.Delete(id)
}
