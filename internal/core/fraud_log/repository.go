package fraud_log

import (
	"github.com/alex/ads_backend/internal/core/fraud_log/dto"
	"gorm.io/gorm"
)

type Repository interface {
	Create(log *FraudLog) error
	Update(log *FraudLog) error
	FindByID(id uint64) (*FraudLogWithNames, error)
	FindAll(filter dto.FraudLogFilter) ([]FraudLogWithNames, int64, error)
	ExistsOpenDuplicate(creativeID, eventType, newValue string) (bool, error)
}

type FraudLogWithNames struct {
	FraudLog
	BrandName             *string `gorm:"column:brand_name"`
	BrandPhoto            *string `gorm:"column:brand_photo"`
	AdAccountName         *string `gorm:"column:ad_account_name"`
	AdAccountBusinessName *string `gorm:"column:ad_account_business_name"`
	CampaignName          *string `gorm:"column:campaign_name"`
	AdSetName             *string `gorm:"column:adset_name"`
	AdName                *string `gorm:"column:ad_name"`
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

func (r *repository) FindByID(id uint64) (*FraudLogWithNames, error) {
	var log FraudLogWithNames
	err := r.db.Model(&FraudLog{}).
		Select("fraud_logs.*, b.name as brand_name, b.photo as brand_photo, acc.name as ad_account_name, acc.business_name as ad_account_business_name, c.name as campaign_name, s.name as adset_name, a.name as ad_name").
		Joins("LEFT JOIN brands b ON fraud_logs.brand_id = b.id").
		Joins("LEFT JOIN meta_ad_accounts acc ON fraud_logs.ad_account_id = acc.id").
		Joins("LEFT JOIN meta_campaigns c ON fraud_logs.campaign_id = c.id").
		Joins("LEFT JOIN meta_ad_sets s ON fraud_logs.adset_id = s.id").
		Joins("LEFT JOIN meta_ads a ON fraud_logs.ad_id = a.id").
		First(&log, id).Error
	return &log, err
}

func (r *repository) FindAll(filter dto.FraudLogFilter) ([]FraudLogWithNames, int64, error) {
	var logs []FraudLogWithNames
	var total int64

	q := r.db.Model(&FraudLog{})
	if filter.BrandID != nil {
		q = q.Where("fraud_logs.brand_id = ?", *filter.BrandID)
	}
	if filter.AdAccountID != "" {
		q = q.Where("fraud_logs.ad_account_id = ?", filter.AdAccountID)
	}
	if filter.CampaignID != "" {
		q = q.Where("fraud_logs.campaign_id = ?", filter.CampaignID)
	}
	if filter.CreativeID != "" {
		q = q.Where("fraud_logs.creative_id = ?", filter.CreativeID)
	}
	if filter.Severity != "" {
		q = q.Where("fraud_logs.severity = ?", filter.Severity)
	}
	if filter.Status != "" {
		q = q.Where("fraud_logs.status = ?", filter.Status)
	}
	if filter.DateStart != "" {
		q = q.Where("fraud_logs.created_at >= ?", filter.DateStart)
	}
	if filter.DateStop != "" {
		q = q.Where("fraud_logs.created_at <= ?", filter.DateStop+" 23:59:59")
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

	q = q.Select("fraud_logs.*, b.name as brand_name, b.photo as brand_photo, acc.name as ad_account_name, acc.business_name as ad_account_business_name, c.name as campaign_name, s.name as adset_name, a.name as ad_name").
		Joins("LEFT JOIN brands b ON fraud_logs.brand_id = b.id").
		Joins("LEFT JOIN meta_ad_accounts acc ON fraud_logs.ad_account_id = acc.id").
		Joins("LEFT JOIN meta_campaigns c ON fraud_logs.campaign_id = c.id").
		Joins("LEFT JOIN meta_ad_sets s ON fraud_logs.adset_id = s.id").
		Joins("LEFT JOIN meta_ads a ON fraud_logs.ad_id = a.id")

	err := q.Order("fraud_logs.created_at DESC").Limit(limit).Offset((page - 1) * limit).Find(&logs).Error
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
