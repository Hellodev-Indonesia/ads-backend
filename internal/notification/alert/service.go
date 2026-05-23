package alert

import (
	"errors"

	"github.com/alex/ads_backend/internal/notification/alert/dto"
)

type Service interface {
	Create(input dto.CreateAlertInput) (*Alert, error)
	FindByID(id uint64) (dto.AlertResponse, error)
	FindAll(filter dto.AlertFilter) ([]dto.AlertResponse, int64, error)
	MarkRead(id uint64) (dto.AlertResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(input dto.CreateAlertInput) (*Alert, error) {
	a := &Alert{
		FraudLogID: input.FraudLogID,
		BrandID:    input.BrandID,
		Title:      input.Title,
		Message:    input.Message,
		Severity:   input.Severity,
		IsRead:     false,
	}
	if err := s.repo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *service) FindByID(id uint64) (dto.AlertResponse, error) {
	a, err := s.repo.FindByID(id)
	if err != nil {
		return dto.AlertResponse{}, errors.New("alert not found")
	}
	return toResponse(*a), nil
}

func (s *service) FindAll(filter dto.AlertFilter) ([]dto.AlertResponse, int64, error) {
	alerts, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, 0, err
	}
	var responses []dto.AlertResponse
	for _, a := range alerts {
		responses = append(responses, toResponse(a))
	}
	return responses, total, nil
}

func (s *service) MarkRead(id uint64) (dto.AlertResponse, error) {
	a, err := s.repo.FindByID(id)
	if err != nil {
		return dto.AlertResponse{}, errors.New("alert not found")
	}
	a.IsRead = true
	if err := s.repo.Update(a); err != nil {
		return dto.AlertResponse{}, err
	}
	return toResponse(*a), nil
}

func toResponse(a Alert) dto.AlertResponse {
	return dto.AlertResponse{
		ID:         a.ID,
		FraudLogID: a.FraudLogID,
		BrandID:    a.BrandID,
		Title:      a.Title,
		Message:    a.Message,
		Severity:   a.Severity,
		IsRead:     a.IsRead,
		CreatedAt:  a.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  a.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
