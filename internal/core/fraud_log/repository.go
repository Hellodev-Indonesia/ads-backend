package fraud_log

import (
	"github.com/alex/ads_backend/internal/core/fraud_log/dto"
	"gorm.io/gorm"
)

type Repository interface {
	Create(log *FraudLog) error
	Update(log *FraudLog) error
	FindByID(id uint64) (*FraudLog, error)
	FindAll(filter dto.FraudLogFilter) ([]FraudLog, int64, error)
	ExistsOpenDuplicate(creativeID, eventType, newValue string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(log *FraudLog) error {
	return r.db.Create(log).Error
}

func (r *repository) Update(log *FraudLog) error {
	return r.db.Save(log).Error
}

func (r *repository) FindByID(id uint64) (*FraudLog, error) {
	var log FraudLog
	err := r.db.First(&log, id).Error
	return &log, err
}

func (r *repository) FindAll(filter dto.FraudLogFilter) ([]FraudLog, int64, error) {
	var logs []FraudLog
	var total int64

	q := r.db.Model(&FraudLog{})
	if filter.BrandID != nil {
		q = q.Where("brand_id = ?", *filter.BrandID)
	}
	if filter.AdAccountID != "" {
		q = q.Where("ad_account_id = ?", filter.AdAccountID)
	}
	if filter.CampaignID != "" {
		q = q.Where("campaign_id = ?", filter.CampaignID)
	}
	if filter.CreativeID != "" {
		q = q.Where("creative_id = ?", filter.CreativeID)
	}
	if filter.Severity != "" {
		q = q.Where("severity = ?", filter.Severity)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.DateStart != "" {
		q = q.Where("created_at >= ?", filter.DateStart)
	}
	if filter.DateStop != "" {
		q = q.Where("created_at <= ?", filter.DateStop+" 23:59:59")
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

	err := q.Order("created_at DESC").Limit(limit).Offset((page - 1) * limit).Find(&logs).Error
	return logs, total, err
}

func (r *repository) ExistsOpenDuplicate(creativeID, eventType, newValue string) (bool, error) {
	var count int64
	err := r.db.Model(&FraudLog{}).
		Where("creative_id = ? AND event_type = ? AND new_value = ? AND status = ?",
			creativeID, eventType, newValue, "open").
		Count(&count).Error
	return count > 0, err
}
