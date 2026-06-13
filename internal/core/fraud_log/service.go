package fraud_log

import (
	"errors"
	"time"

	"github.com/alex/ads_backend/internal/core/fraud_log/dto"
)

type Service interface {
	Create(input dto.CreateFraudLogInput) (*FraudLog, error)
	FindByID(id uint64) (dto.FraudLogResponse, error)
	FindAll(filter dto.FraudLogFilter) ([]dto.FraudLogResponse, int64, error)
	Resolve(id uint64, userID uint64) (dto.FraudLogResponse, error)
	ExistsOpenDuplicate(creativeID, eventType, newValue string) (bool, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(input dto.CreateFraudLogInput) (*FraudLog, error) {
	now := time.Now()
	msg := input.Message
	log := &FraudLog{
		BrandID:       input.BrandID,
		AdAccountID:   input.AdAccountID,
		CampaignID:    input.CampaignID,
		AdsetID:       input.AdsetID,
		AdID:          input.AdID,
		CreativeID:    input.CreativeID,
		EventType:     input.EventType,
		ActorID:       input.ActorID,
		ActorName:     input.ActorName,
		Severity:      input.Severity,
		OldValue:      input.OldValue,
		NewValue:      input.NewValue,
		MatchedRuleID: input.MatchedRuleID,
		Message:       &msg,
		Status:        "open",
		DetectedAt:    &now,
	}
	if err := s.repo.Create(log); err != nil {
		return nil, err
	}
	return log, nil
}

func (s *service) FindByID(id uint64) (dto.FraudLogResponse, error) {
	log, err := s.repo.FindByID(id)
	if err != nil {
		return dto.FraudLogResponse{}, errors.New("fraud log not found")
	}
	return toResponse(*log), nil
}

func (s *service) FindAll(filter dto.FraudLogFilter) ([]dto.FraudLogResponse, int64, error) {
	logs, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, 0, err
	}
	var responses []dto.FraudLogResponse
	for _, l := range logs {
		responses = append(responses, toResponse(l))
	}
	return responses, total, nil
}

func (s *service) Resolve(id uint64, userID uint64) (dto.FraudLogResponse, error) {
	log, err := s.repo.FindByID(id)
	if err != nil {
		return dto.FraudLogResponse{}, errors.New("fraud log not found")
	}
	if log.Status == "resolved" {
		return toResponse(*log), nil
	}
	now := time.Now()
	log.Status = "resolved"
	log.ResolvedAt = &now
	log.ResolvedBy = &userID
	if err := s.repo.Update(&log.FraudLog); err != nil {
		return dto.FraudLogResponse{}, err
	}
	// Need to refetch to get the user name
	updatedLog, _ := s.repo.FindByID(id)
	return toResponse(*updatedLog), nil
}

func (s *service) ExistsOpenDuplicate(creativeID, eventType, newValue string) (bool, error) {
	return s.repo.ExistsOpenDuplicate(creativeID, eventType, newValue)
}

func toResponse(l FraudLogWithNames) dto.FraudLogResponse {
	r := dto.FraudLogResponse{
		ID:            l.ID,
		CreativeID:    l.CreativeID,
		EventType:     l.EventType,
		ActorID:       l.ActorID,
		ActorName:     l.ActorName,
		Severity:      l.Severity,
		OldValue:      l.OldValue,
		NewValue:      l.NewValue,
		MatchedRuleID: l.MatchedRuleID,
		Message:       l.Message,
		Status:        l.Status,
		CreatedAt:     l.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     l.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if l.BrandID != nil && l.BrandName != nil {
		r.Brand = &dto.SimpleBrand{
			ID:    *l.BrandID,
			Name:  *l.BrandName,
			Photo: l.BrandPhoto,
		}
	}
	if l.AdAccountID != nil && l.AdAccountName != nil {
		r.AdAccount = &dto.SimpleAdAccount{
			ID:           *l.AdAccountID,
			Name:         *l.AdAccountName,
			BusinessName: l.AdAccountBusinessName,
		}
	}
	if l.CampaignID != nil && l.CampaignName != nil {
		r.Campaign = &dto.SimpleCampaign{
			ID:   *l.CampaignID,
			Name: *l.CampaignName,
		}
	}
	if l.AdsetID != nil && l.AdSetName != nil {
		r.Adset = &dto.SimpleAdSet{
			ID:   *l.AdsetID,
			Name: *l.AdSetName,
		}
	}
	if l.AdID != nil && l.AdName != nil {
		r.Ad = &dto.SimpleAd{
			ID:   *l.AdID,
			Name: *l.AdName,
		}
	}
	if l.DetectedAt != nil {
		s := l.DetectedAt.Format("2006-01-02 15:04:05")
		r.DetectedAt = &s
	}
	if l.ResolvedAt != nil {
		s := l.ResolvedAt.Format("2006-01-02 15:04:05")
		r.ResolvedAt = &s
	}
	if l.ResolvedBy != nil && l.ResolvedByName != nil {
		r.ResolvedBy = &dto.SimpleUser{
			ID:   *l.ResolvedBy,
			Name: *l.ResolvedByName,
		}
	}
	return r
}
