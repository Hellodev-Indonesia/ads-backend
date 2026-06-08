package insight

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InsightFilter struct {
	AccountID  string
	CampaignID string
	AdSetID    string
	AdID       string
	DateStart  string
	DateStop   string
	Page       int
	Limit      int
}

type Repository interface {
	UpsertBatch(insights []MetaInsight) error
	FindCampaignInsights(filter InsightFilter) ([]MetaInsight, int64, error)
	FindAdInsights(filter InsightFilter) ([]MetaInsight, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) UpsertBatch(insights []MetaInsight) error {
	if len(insights) == 0 {
		return nil
	}
	// For insights with auto-increment PK + unique key, we use ON DUPLICATE KEY UPDATE
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "campaign_id"},
			{Name: "adset_id"},
			{Name: "ad_id"},
			{Name: "level"},
			{Name: "date_start"},
			{Name: "date_stop"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"account_name", "account_currency", "campaign_name",
			"adset_name", "ad_name", "objective",
			"impressions", "reach", "clicks", "inline_link_clicks",
			"inline_link_click_ctr", "spend", "cpc", "cpm", "ctr",
			"actions", "action_values", "cost_per_action_type", "synced_at",
		}),
	}).CreateInBatches(insights, 100).Error
}

func (r *repository) findInsights(level string, filter InsightFilter) ([]MetaInsight, int64, error) {
	var insights []MetaInsight
	var total int64

	query := r.db.Model(&MetaInsight{}).Where("level = ?", level)

	if filter.AccountID != "" {
		query = query.Where("account_id = ?", filter.AccountID)
	}
	if filter.CampaignID != "" {
		query = query.Where("campaign_id = ?", filter.CampaignID)
	}
	if filter.AdSetID != "" {
		query = query.Where("adset_id = ?", filter.AdSetID)
	}
	if filter.AdID != "" {
		query = query.Where("ad_id = ?", filter.AdID)
	}
	if filter.DateStart != "" {
		query = query.Where("date_start >= ?", filter.DateStart)
	}
	if filter.DateStop != "" {
		query = query.Where("date_stop <= ?", filter.DateStop)
	}

	query.Count(&total)

	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit

	err := query.Order("date_start DESC").Limit(filter.Limit).Offset(offset).Find(&insights).Error
	return insights, total, err
}

func (r *repository) FindCampaignInsights(filter InsightFilter) ([]MetaInsight, int64, error) {
	return r.findInsights("campaign", filter)
}

func (r *repository) FindAdInsights(filter InsightFilter) ([]MetaInsight, int64, error) {
	return r.findInsights("ad", filter)
}
