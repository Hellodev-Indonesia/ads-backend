package brand_whitelist_rule

import (
	"errors"

	"github.com/alex/ads_backend/internal/core/brand_whitelist_rule/dto"
)

// MatchResult holds the outcome of an IsURLAllowed check.
type MatchResult struct {
	Allowed bool
	Rule    *BrandWhitelistRule
}

type Service interface {
	Create(brandID uint64, req dto.CreateWhitelistRuleRequest) (dto.WhitelistRuleResponse, error)
	Update(brandID, id uint64, req dto.UpdateWhitelistRuleRequest) (dto.WhitelistRuleResponse, error)
	Delete(brandID, id uint64) error
	FindByID(brandID, id uint64) (dto.WhitelistRuleResponse, error)
	FindAll(brandID uint64, filter dto.WhitelistRuleFilter) ([]dto.WhitelistRuleResponse, int64, error)
	IsURLAllowed(brandID uint64, targetURL string, scope string) (MatchResult, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(brandID uint64, req dto.CreateWhitelistRuleRequest) (dto.WhitelistRuleResponse, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	normalized := computeNormalizedValue(req.MatchType, req.Value)
	rule := &BrandWhitelistRule{
		BrandID:         brandID,
		Scope:           req.Scope,
		MatchType:       req.MatchType,
		Value:           req.Value,
		NormalizedValue: normalized,
		AllowSubdomains: req.AllowSubdomains,
		IsActive:        isActive,
		Description:     req.Description,
	}

	if err := s.repo.Create(rule); err != nil {
		return dto.WhitelistRuleResponse{}, err
	}
	return toResponse(*rule), nil
}

func (s *service) Update(brandID, id uint64, req dto.UpdateWhitelistRuleRequest) (dto.WhitelistRuleResponse, error) {
	rule, err := s.repo.FindByID(id)
	if err != nil || rule.BrandID != brandID {
		return dto.WhitelistRuleResponse{}, errors.New("whitelist rule not found")
	}

	if req.Scope != nil {
		rule.Scope = *req.Scope
	}
	if req.MatchType != nil {
		rule.MatchType = *req.MatchType
	}
	if req.Value != nil {
		rule.Value = *req.Value
	}
	if req.AllowSubdomains != nil {
		rule.AllowSubdomains = *req.AllowSubdomains
	}
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}
	if req.Description != nil {
		rule.Description = req.Description
	}
	rule.NormalizedValue = computeNormalizedValue(rule.MatchType, rule.Value)

	if err := s.repo.Update(rule); err != nil {
		return dto.WhitelistRuleResponse{}, err
	}
	return toResponse(*rule), nil
}

func (s *service) Delete(brandID, id uint64) error {
	rule, err := s.repo.FindByID(id)
	if err != nil || rule.BrandID != brandID {
		return errors.New("whitelist rule not found")
	}
	return s.repo.Delete(id)
}

func (s *service) FindByID(brandID, id uint64) (dto.WhitelistRuleResponse, error) {
	rule, err := s.repo.FindByID(id)
	if err != nil || rule.BrandID != brandID {
		return dto.WhitelistRuleResponse{}, errors.New("whitelist rule not found")
	}
	return toResponse(*rule), nil
}

func (s *service) FindAll(brandID uint64, filter dto.WhitelistRuleFilter) ([]dto.WhitelistRuleResponse, int64, error) {
	rules, total, err := s.repo.FindAllByBrandID(brandID, filter)
	if err != nil {
		return nil, 0, err
	}
	var responses []dto.WhitelistRuleResponse
	for _, r := range rules {
		responses = append(responses, toResponse(r))
	}
	return responses, total, nil
}

// IsURLAllowed checks whether targetURL passes any active whitelist rule for the brand and scope.
// Returns the first matching rule or Allowed=false if none match.
func (s *service) IsURLAllowed(brandID uint64, targetURL string, scope string) (MatchResult, error) {
	rules, err := s.repo.FindActiveByBrandIDAndScope(brandID, scope)
	if err != nil {
		return MatchResult{}, err
	}

	for i := range rules {
		matched, err := matchRule(&rules[i], targetURL)
		if err != nil {
			// Skip rules with invalid regex rather than blocking.
			continue
		}
		if matched {
			return MatchResult{Allowed: true, Rule: &rules[i]}, nil
		}
	}
	return MatchResult{Allowed: false}, nil
}

func computeNormalizedValue(matchType, value string) *string {
	var norm string
	switch matchType {
	case "domain":
		norm = NormalizeDomain(value)
	case "exact_url", "path_prefix":
		norm = NormalizeURL(value)
	default:
		return nil
	}
	return &norm
}

func toResponse(r BrandWhitelistRule) dto.WhitelistRuleResponse {
	return dto.WhitelistRuleResponse{
		ID:              r.ID,
		BrandID:         r.BrandID,
		Scope:           r.Scope,
		MatchType:       r.MatchType,
		Value:           r.Value,
		NormalizedValue: r.NormalizedValue,
		AllowSubdomains: r.AllowSubdomains,
		IsActive:        r.IsActive,
		Description:     r.Description,
		CreatedAt:       r.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       r.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
