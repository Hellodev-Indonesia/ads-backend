package brand_whitelist_rule

import (
	"github.com/alex/ads_backend/internal/core/brand_whitelist_rule/dto"
	"gorm.io/gorm"
)

type Repository interface {
	Create(rule *BrandWhitelistRule) error
	Update(rule *BrandWhitelistRule) error
	Delete(id uint64) error
	FindByID(id uint64) (*BrandWhitelistRule, error)
	FindAllByBrandID(brandID uint64, filter dto.WhitelistRuleFilter) ([]BrandWhitelistRule, int64, error)
	FindActiveByBrandIDAndScope(brandID uint64, scope string) ([]BrandWhitelistRule, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(rule *BrandWhitelistRule) error {
	return r.db.Create(rule).Error
}

func (r *repository) Update(rule *BrandWhitelistRule) error {
	return r.db.Save(rule).Error
}

func (r *repository) Delete(id uint64) error {
	return r.db.Delete(&BrandWhitelistRule{}, id).Error
}

func (r *repository) FindByID(id uint64) (*BrandWhitelistRule, error) {
	var rule BrandWhitelistRule
	err := r.db.First(&rule, id).Error
	return &rule, err
}

func (r *repository) FindAllByBrandID(brandID uint64, filter dto.WhitelistRuleFilter) ([]BrandWhitelistRule, int64, error) {
	var rules []BrandWhitelistRule
	var total int64

	q := r.db.Model(&BrandWhitelistRule{}).Where("brand_id = ?", brandID)
	if filter.Scope != "" {
		q = q.Where("scope = ?", filter.Scope)
	}
	if filter.MatchType != "" {
		q = q.Where("match_type = ?", filter.MatchType)
	}
	if filter.IsActive != nil {
		q = q.Where("is_active = ?", *filter.IsActive)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := filter.Limit
	page := filter.Page
	if limit <= 0 {
		limit = 25
	}
	if page <= 0 {
		page = 1
	}

	err := q.Order("created_at DESC").Limit(limit).Offset((page - 1) * limit).Find(&rules).Error
	return rules, total, err
}

func (r *repository) FindActiveByBrandIDAndScope(brandID uint64, scope string) ([]BrandWhitelistRule, error) {
	var rules []BrandWhitelistRule
	err := r.db.Where("brand_id = ? AND scope = ? AND is_active = 1", brandID, scope).Find(&rules).Error
	return rules, err
}
