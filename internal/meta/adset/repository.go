package adset

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/alex/ads_backend/internal/meta/adset/dto"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AdSetFilter struct {
	AccountID   string
	BrandID     *uint64
	CampaignID  string
	CampaignIDs []string
	AdSetIDs    []string
	Status      string
	Search      string
	DateStart  string
	DateStop   string
	Page       int
	Limit      int
}

type Repository interface {
	Upsert(adset *MetaAdSet) error
	UpsertBatch(adsets []MetaAdSet) error
	FindAll(filter AdSetFilter) ([]MetaAdSet, int64, error)
	FindByCampaignID(campaignID string) ([]MetaAdSet, error)
	FindAdSetDashboard(filter AdSetFilter) ([]adSetDashboardScan, int64, error)
	FindSimpleListByBrand(brandID uint64, campaignIDs []string) ([]dto.SimpleListResponse, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Upsert(adset *MetaAdSet) error {
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(adset).Error
}

func (r *repository) UpsertBatch(adsets []MetaAdSet) error {
	if len(adsets) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(adsets, 100).Error
}

func (r *repository) FindAll(filter AdSetFilter) ([]MetaAdSet, int64, error) {
	var adsets []MetaAdSet
	var total int64

	query := r.db.Model(&MetaAdSet{})

	if filter.CampaignID != "" {
		query = query.Where("campaign_id = ?", filter.CampaignID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		query = query.Where("name LIKE ?", "%"+filter.Search+"%")
	}
	if filter.BrandID != nil {
		query = query.Where("campaign_id IN (SELECT id FROM meta_campaigns WHERE account_id IN (SELECT id FROM meta_ad_accounts WHERE brand_id = ?))", *filter.BrandID)
	}

	query.Count(&total)

	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit

	err := query.Order("created_time DESC").Limit(filter.Limit).Offset(offset).Find(&adsets).Error
	return adsets, total, err
}

func (r *repository) FindByCampaignID(campaignID string) ([]MetaAdSet, error) {
	var adsets []MetaAdSet
	err := r.db.Where("campaign_id = ?", campaignID).Find(&adsets).Error
	return adsets, err
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
	BidStrategy       string          `gorm:"column:bid_strategy"`
	UpdatedTime       *time.Time      `gorm:"column:updated_time"`
	StartTime         *time.Time      `gorm:"column:start_time"`
}
func (r *repository) FindAdSetDashboard(filter AdSetFilter) ([]adSetDashboardScan, int64, error) {
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
	if len(filter.CampaignIDs) > 0 {
		where += " AND s.campaign_id IN ?"
		args = append(args, filter.CampaignIDs)
	}
	if len(filter.AdSetIDs) > 0 {
		where += " AND s.id IN ?"
		args = append(args, filter.AdSetIDs)
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
  s.start_time,
  s.end_time,
  s.updated_time,
  s.bid_strategy,
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

func javaStrToInt(s string) (int, error) {
	var v int
	_, _ = fmt.Sscanf(s, "%d", &v)
	return v, nil
}

func intToStr(v int) string {
	return fmt.Sprintf("%d", v)
}

func (r *repository) FindSimpleListByBrand(brandID uint64, campaignIDs []string) ([]dto.SimpleListResponse, error) {
	var list []dto.SimpleListResponse
	query := r.db.Table("meta_ad_sets s").
		Select("s.id, s.name").
		Joins("JOIN meta_campaigns c ON s.campaign_id = c.id").
		Joins("JOIN meta_ad_accounts a ON c.account_id = a.id").
		Where("a.brand_id = ?", brandID).
		Order("s.created_time DESC")

	if len(campaignIDs) > 0 {
		query = query.Where("s.campaign_id IN ?", campaignIDs)
	}

	err := query.Scan(&list).Error
	return list, err
}
