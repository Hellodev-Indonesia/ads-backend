package dashboard

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetSummaryMetrics(ctx context.Context) (float64, int64, int64, error)
	GetBrandsMonitoring(ctx context.Context) ([]BrandMonitoringItemRaw, error)
	GetRecentActivities(ctx context.Context) ([]RecentActivityItemRaw, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

type BrandMonitoringItemRaw struct {
	BrandID        uint64
	BrandName      string
	AdAccountCount int64
	ActiveCampaign int64
	TotalSpends    float64
}

type RecentActivityItemRaw struct {
	ActivityID      uint64
	BrandName       string
	AdAccountName   string
	Activity        string
	ActivityDetails string
	ItemChanged     string
	ChangeBy        string
	EventTime       *time.Time
}

func (r *dashboardRepository) GetSummaryMetrics(ctx context.Context) (float64, int64, int64, error) {
	var totalSpend float64
	var totalOngoingAds int64
	var securityIssues int64

	// Total Spend
	if err := r.db.WithContext(ctx).Table("meta_insights").Select("COALESCE(SUM(spend), 0)").Scan(&totalSpend).Error; err != nil {
		return 0, 0, 0, err
	}

	// Total Ongoing Ads
	if err := r.db.WithContext(ctx).Table("meta_ads").Where("effective_status = ?", "ACTIVE").Count(&totalOngoingAds).Error; err != nil {
		return 0, 0, 0, err
	}

	// Security Issues
	if err := r.db.WithContext(ctx).Table("fraud_logs").Where("status = ?", "open").Count(&securityIssues).Error; err != nil {
		return 0, 0, 0, err
	}

	return totalSpend, totalOngoingAds, securityIssues, nil
}

func (r *dashboardRepository) GetBrandsMonitoring(ctx context.Context) ([]BrandMonitoringItemRaw, error) {
	var items []BrandMonitoringItemRaw

	query := `
		SELECT 
			b.id as brand_id,
			b.name as brand_name,
			COUNT(DISTINCT a.id) as ad_account_count,
			COUNT(DISTINCT CASE WHEN c.effective_status = 'ACTIVE' THEN c.id END) as active_campaign,
			COALESCE(SUM(i.spend), 0) as total_spends
		FROM brands b
		LEFT JOIN meta_ad_accounts a ON b.id = a.brand_id
		LEFT JOIN meta_campaigns c ON a.id = c.account_id
		LEFT JOIN meta_insights i ON c.id = i.campaign_id
		GROUP BY b.id, b.name
	`

	if err := r.db.WithContext(ctx).Raw(query).Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *dashboardRepository) GetRecentActivities(ctx context.Context) ([]RecentActivityItemRaw, error) {
	var items []RecentActivityItemRaw

	query := `
		SELECT 
			act.id as activity_id,
			COALESCE(b.name, 'Unknown') as brand_name,
			COALESCE(a.name, act.ad_account_id) as ad_account_name,
			COALESCE(act.event_type, 'Activity') as activity,
			COALESCE(act.object_name, 'Unknown') as activity_details,
			COALESCE(act.object_type, 'Item') as item_changed,
			COALESCE(act.actor_name, 'System') as change_by,
			act.event_time as event_time
		FROM meta_activities act
		LEFT JOIN meta_ad_accounts a ON act.ad_account_id = a.id
		LEFT JOIN brands b ON a.brand_id = b.id
		ORDER BY act.event_time DESC
		LIMIT 50
	`

	if err := r.db.WithContext(ctx).Raw(query).Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
