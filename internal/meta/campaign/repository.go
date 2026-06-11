package campaign

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CampaignFilter struct {
	Status    string
	AccountID string
	BrandID   *uint64
	Search    string
	DateStart string
	DateStop  string
	Page      int
	Limit     int
}

type InsightSummaryRow struct {
	Spend       float64
	Impressions int64
	Reach       int64
	Actions     []byte // JSON
}

type Repository interface {
	Upsert(campaign *MetaCampaign) error
	UpsertBatch(campaigns []MetaCampaign) error
	FindAll(filter CampaignFilter) ([]MetaCampaign, int64, error)
	FindByID(id string) (*MetaCampaign, error)
	GetSummaryByBrand(brandID uint64, dateStart, dateStop string) ([]InsightSummaryRow, error)
	FindCampaignDashboard(filter CampaignFilter) ([]campaignDashboardScan, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Upsert(campaign *MetaCampaign) error {
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(campaign).Error
}

func (r *repository) UpsertBatch(campaigns []MetaCampaign) error {
	if len(campaigns) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(campaigns, 100).Error
}

func (r *repository) FindAll(filter CampaignFilter) ([]MetaCampaign, int64, error) {
	var campaigns []MetaCampaign
	var total int64

	query := r.db.Model(&MetaCampaign{})

	if filter.AccountID != "" {
		query = query.Where("account_id = ?", filter.AccountID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		query = query.Where("name LIKE ?", "%"+filter.Search+"%")
	}
	if filter.BrandID != nil {
		query = query.Where("account_id IN (SELECT id FROM meta_ad_accounts WHERE brand_id = ?)", *filter.BrandID)
	}

	query.Count(&total)

	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit

	err := query.Order("created_time DESC").Limit(filter.Limit).Offset(offset).Find(&campaigns).Error
	return campaigns, total, err
}

func (r *repository) FindByID(id string) (*MetaCampaign, error) {
	var campaign MetaCampaign
	err := r.db.Where("id = ?", id).First(&campaign).Error
	return &campaign, err
}

func (r *repository) GetSummaryByBrand(brandID uint64, dateStart, dateStop string) ([]InsightSummaryRow, error) {
	var rows []InsightSummaryRow
	query := r.db.Table("meta_insights").
		Select("spend, impressions, reach, actions").
		Where("level = 'campaign'").
		Where("campaign_id IN (SELECT id FROM meta_campaigns WHERE account_id IN (SELECT id FROM meta_ad_accounts WHERE brand_id = ?))", brandID)

	if dateStart != "" {
		query = query.Where("date_start >= ?", dateStart)
	}
	if dateStop != "" {
		query = query.Where("date_stop <= ?", dateStop)
	}

	err := query.Scan(&rows).Error
	return rows, err
}

type campaignDashboardScan struct {
	CampaignID        string          `gorm:"column:campaign_id"`
	CampaignName      string          `gorm:"column:campaign_name"`
	Status            string          `gorm:"column:status"`
	EffectiveStatus   string          `gorm:"column:effective_status"`
	Objective         string          `gorm:"column:objective"`
	DailyBudget         float64         `gorm:"column:daily_budget"`
	LifetimeBudget      float64         `gorm:"column:lifetime_budget"`
	AdsetDailyBudget    float64         `gorm:"column:adset_daily_budget"`
	AdsetLifetimeBudget float64         `gorm:"column:adset_lifetime_budget"`
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
func (r *repository) FindCampaignDashboard(filter CampaignFilter) ([]campaignDashboardScan, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	where := "WHERE 1=1"
	var args []interface{}

	if filter.AccountID != "" {
		where += " AND c.account_id = ?"
		args = append(args, filter.AccountID)
	}

	if filter.BrandID != nil {
		where += " AND c.account_id IN (SELECT id FROM meta_ad_accounts WHERE brand_id = ?)"
		args = append(args, *filter.BrandID)
	}

	if filter.Status != "" {
		where += " AND c.status = ?"
		args = append(args, filter.Status)
	}
	if filter.Search != "" {
		where += " AND c.name LIKE ?"
		args = append(args, "%"+filter.Search+"%")
	}

	var total int64
	countSQL := "SELECT COUNT(*) FROM meta_campaigns c " + where
	if err := r.db.Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit

	// Fetch campaigns and attribution spec
	campaignsSQL := `
SELECT
  c.id            AS campaign_id,
  c.name          AS campaign_name,
  c.status,
  c.effective_status,
  c.objective,
  c.daily_budget,
  c.lifetime_budget,
  c.stop_time,
  a.attribution_spec,
  COALESCE(s_sum.sum_daily, 0) AS adset_daily_budget,
  COALESCE(s_sum.sum_lifetime, 0) AS adset_lifetime_budget
FROM meta_campaigns c
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
LEFT JOIN (
  SELECT campaign_id, SUM(daily_budget) as sum_daily, SUM(lifetime_budget) as sum_lifetime
  FROM meta_ad_sets
  GROUP BY campaign_id
) s_sum ON c.id = s_sum.campaign_id
` + where + ` ORDER BY c.created_time DESC LIMIT ? OFFSET ?`

	queryArgs := append(args, filter.Limit, offset)

	var rows []campaignDashboardScan
	if err := r.db.Raw(campaignsSQL, queryArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return rows, total, nil
	}

	// Gather campaign IDs to fetch insights
	var campaignIDs []string
	for _, r := range rows {
		campaignIDs = append(campaignIDs, r.CampaignID)
	}

	// Fetch insights
	insightWhere := "level = 'campaign' AND campaign_id IN ?"
	insightArgs := []interface{}{campaignIDs}

	if filter.DateStart != "" {
		insightWhere += " AND date_start >= ?"
		insightArgs = append(insightArgs, filter.DateStart)
	}
	if filter.DateStop != "" {
		insightWhere += " AND date_stop <= ?"
		insightArgs = append(insightArgs, filter.DateStop)
	}

	type rawInsight struct {
		CampaignID        string
		Spend             float64
		Impressions       int64
		Reach             int64
		Actions           string
		CostPerActionType string
		DateStart         string
		DateStop          string
	}
	var rawInsights []rawInsight
	if err := r.db.Table("meta_insights").Select("campaign_id, spend, impressions, reach, CAST(actions AS CHAR) as actions, CAST(cost_per_action_type AS CHAR) as cost_per_action_type, date_start, date_stop").Where(insightWhere, insightArgs...).Find(&rawInsights).Error; err != nil {
		return nil, 0, err
	}

	// Aggregate in Go
	insightMap := make(map[string]*campaignDashboardScan)
	for i := range rows {
		insightMap[rows[i].CampaignID] = &rows[i]
	}

	type actionItem struct {
		ActionType string `json:"action_type"`
		Value      string `json:"value"`
	}

	for _, ins := range rawInsights {
		row := insightMap[ins.CampaignID]
		if row == nil {
			continue
		}

		if row.Spend == nil {
			val := ins.Spend
			row.Spend = &val
		} else {
			*row.Spend += ins.Spend
		}

		if row.Impressions == nil {
			val := ins.Impressions
			row.Impressions = &val
		} else {
			*row.Impressions += ins.Impressions
		}

		// Reach isn't strictly summable, but for dashboard approximation we sum it
		if row.Reach == nil {
			val := ins.Reach
			row.Reach = &val
		} else {
			*row.Reach += ins.Reach
		}

		if row.DateStart == nil || ins.DateStart < *row.DateStart {
			ds := ins.DateStart
			row.DateStart = &ds
		}
		if row.DateStop == nil || ins.DateStop > *row.DateStop {
			ds := ins.DateStop
			row.DateStop = &ds
		}

		// Merge Actions
		if ins.Actions != "" && ins.Actions != "null" {
			var current []actionItem
			if row.Actions != nil {
				_ = json.Unmarshal(row.Actions, &current)
			}
			var newActions []actionItem
			_ = json.Unmarshal([]byte(ins.Actions), &newActions)

			// sum them
			for _, na := range newActions {
				found := false
				for j, ca := range current {
					if ca.ActionType == na.ActionType {
						cv, _ := javaStrToInt(ca.Value)
						nv, _ := javaStrToInt(na.Value)
						current[j].Value = intToStr(cv + nv)
						found = true
						break
					}
				}
				if !found {
					current = append(current, na)
				}
			}
			b, _ := json.Marshal(current)
			row.Actions = b
		}
	}

	return rows, total, nil
}

func javaStrToInt(s string) (int, error) {
	var v int
	_, _ = fmt.Sscanf(s, "%d", &v)
	return v, nil
}

func intToStr(v int) string {
	return fmt.Sprintf("%d", v)
}