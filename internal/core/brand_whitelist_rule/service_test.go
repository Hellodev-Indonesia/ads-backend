package brand_whitelist_rule_test

import (
	"testing"
	"time"

	"github.com/alex/ads_backend/internal/core/brand_whitelist_rule"
	"github.com/alex/ads_backend/internal/core/brand_whitelist_rule/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupService(t *testing.T) (*brand_whitelist_rule.MockRepository, brand_whitelist_rule.Service) {
	mockRepo := brand_whitelist_rule.NewMockRepository(t)
	svc := brand_whitelist_rule.NewService(mockRepo)
	return mockRepo, svc
}

func TestService_Create(t *testing.T) {
	mockRepo, svc := setupService(t)

	isActive := true
	req := dto.CreateWhitelistRuleRequest{
		Scope:           "domain",
		MatchType:       "domain",
		Value:           "example.com",
		AllowSubdomains: true,
		IsActive:        &isActive,
	}

	mockRepo.On("Create", mock.AnythingOfType("*brand_whitelist_rule.BrandWhitelistRule")).Run(func(args mock.Arguments) {
		r := args.Get(0).(*brand_whitelist_rule.BrandWhitelistRule)
		r.ID = 1
		r.BrandID = 1
		r.CreatedAt = time.Now()
		r.UpdatedAt = time.Now()
	}).Return(nil)

	resp, err := svc.Create(1, req)

	assert.NoError(t, err)
	assert.Equal(t, uint64(1), resp.ID)
	assert.Equal(t, uint64(1), resp.BrandID)
	assert.Equal(t, "example.com", resp.Value)
}

func TestService_Update(t *testing.T) {
	mockRepo, svc := setupService(t)

	newValue := "newdomain.com"
	req := dto.UpdateWhitelistRuleRequest{
		Value: &newValue,
	}

	existing := &brand_whitelist_rule.BrandWhitelistRule{
		ID:        1,
		BrandID:   1,
		Value:     "old.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByID", uint64(1)).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*brand_whitelist_rule.BrandWhitelistRule")).Return(nil)

	resp, err := svc.Update(1, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, "newdomain.com", resp.Value)
}

func TestService_Delete(t *testing.T) {
	mockRepo, svc := setupService(t)

	existing := &brand_whitelist_rule.BrandWhitelistRule{ID: 1, BrandID: 1}
	mockRepo.On("FindByID", uint64(1)).Return(existing, nil)
	mockRepo.On("Delete", uint64(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
}

func TestService_FindByID(t *testing.T) {
	mockRepo, svc := setupService(t)

	existing := &brand_whitelist_rule.BrandWhitelistRule{
		ID:        1,
		BrandID:   1,
		Value:     "example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByID", uint64(1)).Return(existing, nil)

	resp, err := svc.FindByID(1, 1)

	assert.NoError(t, err)
	assert.Equal(t, uint64(1), resp.ID)
	assert.Equal(t, "example.com", resp.Value)
}

func TestService_FindAll(t *testing.T) {
	mockRepo, svc := setupService(t)

	filter := dto.WhitelistRuleFilter{
		Page:  1,
		Limit: 10,
	}

	rules := []brand_whitelist_rule.BrandWhitelistRule{
		{ID: 1, BrandID: 1, Value: "example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	mockRepo.On("FindAllByBrandID", uint64(1), filter).Return(rules, int64(1), nil)

	resp, total, err := svc.FindAll(1, filter)

	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, int64(1), total)
}

func TestService_IsURLAllowed(t *testing.T) {
	mockRepo, svc := setupService(t)

	rules := []brand_whitelist_rule.BrandWhitelistRule{
		{ID: 1, BrandID: 1, MatchType: "domain", Value: "example.com", NormalizedValue: func() *string { s := "example.com"; return &s }(), AllowSubdomains: true},
	}

	mockRepo.On("FindActiveByBrandIDAndScope", uint64(1), "domain").Return(rules, nil)

	result, err := svc.IsURLAllowed(1, "https://sub.example.com", "domain")

	assert.NoError(t, err)
	assert.True(t, result.Allowed)
	assert.NotNil(t, result.Rule)
	assert.Equal(t, uint64(1), result.Rule.ID)
}

func TestService_IsURLAllowed_Blocked(t *testing.T) {
	mockRepo, svc := setupService(t)

	rules := []brand_whitelist_rule.BrandWhitelistRule{
		{ID: 1, BrandID: 1, MatchType: "domain", Value: "example.com", NormalizedValue: func() *string { s := "example.com"; return &s }(), AllowSubdomains: false},
	}

	mockRepo.On("FindActiveByBrandIDAndScope", uint64(1), "domain").Return(rules, nil)

	result, err := svc.IsURLAllowed(1, "https://sub.example.com", "domain")

	assert.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Nil(t, result.Rule)
}
