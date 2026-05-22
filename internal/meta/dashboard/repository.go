package dashboard

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type DashboardFilter struct {
	AccountID  string
	BrandID    *uint64
	CampaignID string
	AdSetID    string
	Status     string
	Search     string
	DateStart  string
	DateStop   string
	Page       int
	Limit      int
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

type adSetDashboardScan struct {
	AdSetID           string          `gorm:"column:adset_id"`
	CampaignID        string          `gorm:"column:campaign_id"`
	CampaignName      string          `gorm:"column:campaign_name"`
	AdSetName         string          `gorm:"column:adset_name"`
	Status            string          `gorm:"column:status"`
	EffectiveStatus   string          `gorm:"column:effective_status"`
	DailyBudget       float64         `gorm:"column:daily_budget"`
	LifetimeBudget    float64         `gorm:"column:lifetime_budget"`
	EndTime           *time.Time      `gorm:"column:end_time"`
	AttributionSpec   json.RawMessage `gorm:"column:attribution_spec"`
	Spend             *float64        `gorm:"column:spend"`
	Impressions       *int64          `gorm:"column:impressions"`
	Reach             *int64          `gorm:"column:reach"`
	Actions           json.RawMessage `gorm:"column:actions"`
	CostPerActionType json.RawMessage `gorm:"column:cost_per_action_type"`
	DateStart         *string         `gorm:"column:date_start"`
	DateStop          *string         `gorm:"column:date_stop"`
}

type adDashboardScan struct {
	AdID              string          `gorm:"column:ad_id"`
	AdSetID           string          `gorm:"column:adset_id"`
	CampaignID        string          `gorm:"column:campaign_id"`
	CampaignName      string          `gorm:"column:campaign_name"`
	AdSetName         string          `gorm:"column:adset_name"`
	AdName            string          `gorm:"column:ad_name"`
	Status            string          `gorm:"column:status"`
	EffectiveStatus   string          `gorm:"column:effective_status"`
	CreativeID        string          `gorm:"column:creative_id"`
	Spend             *float64        `gorm:"column:spend"`
	Impressions       *int64          `gorm:"column:impressions"`
	Reach             *int64          `gorm:"column:reach"`
	Actions           json.RawMessage `gorm:"column:actions"`
	CostPerActionType json.RawMessage `gorm:"column:cost_per_action_type"`
	DateStart         *string         `gorm:"column:date_start"`
	DateStop          *string         `gorm:"column:date_stop"`
}

type Repository interface {
	FindCampaignDashboard(filter DashboardFilter) ([]campaignDashboardScan, int64, error)
	FindAdSetDashboard(filter DashboardFilter) ([]adSetDashboardScan, int64, error)
	FindAdDashboard(filter DashboardFilter) ([]adDashboardScan, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindCampaignDashboard(filter DashboardFilter) ([]campaignDashboardScan, int64, error) {
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
  a.attribution_spec
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
	importStrconv := true
	_ = importStrconv
	var v int
	_, _ = fmt.Sscanf(s, "%d", &v)
	return v, nil
}

func intToStr(v int) string {
	return fmt.Sprintf("%d", v)
}

func (r *repository) FindAdSetDashboard(filter DashboardFilter) ([]adSetDashboardScan, int64, error) {
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
	if filter.CampaignID != "" {
		where += " AND s.campaign_id = ?"
		args = append(args, filter.CampaignID)
	}
	if filter.Status != "" {
		where += " AND s.status = ?"
		args = append(args, filter.Status)
	}
	if filter.Search != "" {
		where += " AND s.name LIKE ?"
		args = append(args, "%"+filter.Search+"%")
	}

	var total int64
	countSQL := "SELECT COUNT(*) FROM meta_ad_sets s JOIN meta_campaigns c ON s.campaign_id = c.id " + where
	if err := r.db.Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit

	adsetsSQL := `
SELECT
  s.id               AS adset_id,
  c.id               AS campaign_id,
  c.name             AS campaign_name,
  s.name             AS adset_name,
  s.status,
  s.effective_status,
  s.daily_budget,
  s.lifetime_budget,
  s.end_time,
  s.attribution_spec
FROM meta_ad_sets s
JOIN meta_campaigns c ON s.campaign_id = c.id
` + where + ` ORDER BY s.created_time DESC LIMIT ? OFFSET ?`

	queryArgs := append(args, filter.Limit, offset)

	var rows []adSetDashboardScan
	if err := r.db.Raw(adsetsSQL, queryArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return rows, total, nil
	}

	var adSetIDs []string
	for _, r := range rows {
		adSetIDs = append(adSetIDs, r.AdSetID)
	}

	insightWhere := "level = 'ad' AND adset_id IN ?"
	insightArgs := []interface{}{adSetIDs}

	if filter.DateStart != "" {
		insightWhere += " AND date_start >= ?"
		insightArgs = append(insightArgs, filter.DateStart)
	}
	if filter.DateStop != "" {
		insightWhere += " AND date_stop <= ?"
		insightArgs = append(insightArgs, filter.DateStop)
	}

	type rawInsight struct {
		AdSetID           string `gorm:"column:adset_id"`
		Spend             float64
		Impressions       int64
		Reach             int64
		Actions           string
		CostPerActionType string
		DateStart         string
		DateStop          string
	}
	var rawInsights []rawInsight
	if err := r.db.Table("meta_insights").Select("adset_id, spend, impressions, reach, CAST(actions AS CHAR) as actions, CAST(cost_per_action_type AS CHAR) as cost_per_action_type, date_start, date_stop").Where(insightWhere, insightArgs...).Find(&rawInsights).Error; err != nil {
		return nil, 0, err
	}

	insightMap := make(map[string]*adSetDashboardScan)
	for i := range rows {
		insightMap[rows[i].AdSetID] = &rows[i]
	}

	type actionItem struct {
		ActionType string `json:"action_type"`
		Value      string `json:"value"`
	}

	for _, ins := range rawInsights {
		row := insightMap[ins.AdSetID]
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

		if ins.Actions != "" && ins.Actions != "null" {
			var current []actionItem
			if row.Actions != nil {
				_ = json.Unmarshal(row.Actions, &current)
			}
			var newActions []actionItem
			_ = json.Unmarshal([]byte(ins.Actions), &newActions)

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

func (r *repository) FindAdDashboard(filter DashboardFilter) ([]adDashboardScan, int64, error) {
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
	if filter.CampaignID != "" {
		where += " AND a.campaign_id = ?"
		args = append(args, filter.CampaignID)
	}
	if filter.AdSetID != "" {
		where += " AND a.adset_id = ?"
		args = append(args, filter.AdSetID)
	}
	if filter.Status != "" {
		where += " AND a.status = ?"
		args = append(args, filter.Status)
	}
	if filter.Search != "" {
		where += " AND a.name LIKE ?"
		args = append(args, "%"+filter.Search+"%")
	}

	var total int64
	countSQL := "SELECT COUNT(*) FROM meta_ads a JOIN meta_campaigns c ON a.campaign_id = c.id " + where
	if err := r.db.Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit

	adsSQL := `
SELECT
  a.id               AS ad_id,
  a.adset_id         AS adset_id,
  s.name             AS adset_name,
  c.id               AS campaign_id,
  c.name             AS campaign_name,
  a.name             AS ad_name,
  a.status,
  a.effective_status,
  a.creative_id
FROM meta_ads a
JOIN meta_ad_sets s ON a.adset_id = s.id
JOIN meta_campaigns c ON a.campaign_id = c.id
` + where + ` ORDER BY a.created_time DESC LIMIT ? OFFSET ?`

	queryArgs := append(args, filter.Limit, offset)

	var rows []adDashboardScan
	if err := r.db.Raw(adsSQL, queryArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return rows, total, nil
	}

	var adIDs []string
	for _, r := range rows {
		adIDs = append(adIDs, r.AdID)
	}

	insightWhere := "level = 'ad' AND ad_id IN ?"
	insightArgs := []interface{}{adIDs}

	if filter.DateStart != "" {
		insightWhere += " AND date_start >= ?"
		insightArgs = append(insightArgs, filter.DateStart)
	}
	if filter.DateStop != "" {
		insightWhere += " AND date_stop <= ?"
		insightArgs = append(insightArgs, filter.DateStop)
	}

	type rawInsight struct {
		AdID              string `gorm:"column:ad_id"`
		Spend             float64
		Impressions       int64
		Reach             int64
		Actions           string
		CostPerActionType string
		DateStart         string
		DateStop          string
	}
	var rawInsights []rawInsight
	if err := r.db.Table("meta_insights").Select("ad_id, spend, impressions, reach, CAST(actions AS CHAR) as actions, CAST(cost_per_action_type AS CHAR) as cost_per_action_type, date_start, date_stop").Where(insightWhere, insightArgs...).Find(&rawInsights).Error; err != nil {
		return nil, 0, err
	}

	insightMap := make(map[string]*adDashboardScan)
	for i := range rows {
		insightMap[rows[i].AdID] = &rows[i]
	}

	type actionItem struct {
		ActionType string `json:"action_type"`
		Value      string `json:"value"`
	}

	for _, ins := range rawInsights {
		row := insightMap[ins.AdID]
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

		if ins.Actions != "" && ins.Actions != "null" {
			var current []actionItem
			if row.Actions != nil {
				_ = json.Unmarshal(row.Actions, &current)
			}
			var newActions []actionItem
			_ = json.Unmarshal([]byte(ins.Actions), &newActions)

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
