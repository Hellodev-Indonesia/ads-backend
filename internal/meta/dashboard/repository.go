package dashboard

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type DashboardFilter struct {
	AccountID string
	Status    string
	Search    string
	Page      int
	Limit     int
}

// campaignDashboardScan holds the raw SQL scan result before mapping to DTO.
type campaignDashboardScan struct {
	CampaignID        string          `gorm:"column:campaign_id"`
	CampaignName      string          `gorm:"column:campaign_name"`
	Status            string          `gorm:"column:status"`
	EffectiveStatus   string          `gorm:"column:effective_status"`
	Objective         string          `gorm:"column:objective"`
	DailyBudget       float64         `gorm:"column:daily_budget"`
	LifetimeBudget    float64         `gorm:"column:lifetime_budget"`
	StopTime          *time.Time      `gorm:"column:stop_time"`
	Spend             *float64        `gorm:"column:spend"`
	Impressions       *int64          `gorm:"column:impressions"`
	Reach             *int64          `gorm:"column:reach"`
	Actions           json.RawMessage `gorm:"column:actions"`
	CostPerActionType json.RawMessage `gorm:"column:cost_per_action_type"`
	AttributionSpec   json.RawMessage `gorm:"column:attribution_spec"`
	DateStart         *string         `gorm:"column:date_start"`
	DateStop          *string         `gorm:"column:date_stop"`
}

type Repository interface {
	FindCampaignDashboard(filter DashboardFilter) ([]campaignDashboardScan, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

const campaignDashboardSQL = `
SELECT
  c.id            AS campaign_id,
  c.name          AS campaign_name,
  c.status,
  c.effective_status,
  c.objective,
  c.daily_budget,
  c.lifetime_budget,
  c.stop_time,
  i.spend,
  i.impressions,
  i.reach,
  i.actions,
  i.cost_per_action_type,
  i.date_start,
  i.date_stop,
  a.attribution_spec
FROM meta_campaigns c
LEFT JOIN (
  SELECT i1.*
  FROM meta_insights i1
  INNER JOIN (
    SELECT campaign_id, MAX(date_start) AS max_date
    FROM meta_insights
    WHERE level = 'campaign'
    GROUP BY campaign_id
  ) latest ON i1.campaign_id = latest.campaign_id
         AND i1.date_start   = latest.max_date
         AND i1.level        = 'campaign'
) i ON c.id = i.campaign_id
LEFT JOIN (
  SELECT a1.campaign_id, a1.attribution_spec
  FROM meta_ad_sets a1
  INNER JOIN (
    SELECT campaign_id, MIN(id) AS min_id
    FROM meta_ad_sets
    GROUP BY campaign_id
  ) first_ad ON a1.campaign_id = first_ad.campaign_id
           AND a1.id           = first_ad.min_id
) a ON c.id = a.campaign_id
`

func (r *repository) FindCampaignDashboard(filter DashboardFilter) ([]campaignDashboardScan, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	where := "WHERE c.account_id = ?"
	args := []interface{}{filter.AccountID}

	if filter.Status != "" {
		where += " AND c.status = ?"
		args = append(args, filter.Status)
	}
	if filter.Search != "" {
		where += " AND c.name LIKE ?"
		args = append(args, "%"+filter.Search+"%")
	}

	// Count total matching campaigns
	var total int64
	countSQL := "SELECT COUNT(*) FROM meta_campaigns c " + where
	if err := r.db.Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	dataSQL := campaignDashboardSQL + where + " ORDER BY c.created_time DESC LIMIT ? OFFSET ?"
	queryArgs := append(args, filter.Limit, offset)

	var rows []campaignDashboardScan
	if err := r.db.Raw(dataSQL, queryArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}
