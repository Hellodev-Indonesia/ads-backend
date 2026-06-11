package ads

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AdFilter struct {
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

type Repository interface {
	Upsert(ad *MetaAd) error
	UpsertBatch(ads []MetaAd) error
	FindAll(filter AdFilter) ([]MetaAd, int64, error)
	FindCreativeRawPayload(creativeID string) (string, error)
	FindAdDashboard(filter AdFilter) ([]adDashboardScan, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Upsert(ad *MetaAd) error {
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(ad).Error
}

func (r *repository) UpsertBatch(ads []MetaAd) error {
	if len(ads) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(ads, 100).Error
}

func (r *repository) FindAll(filter AdFilter) ([]MetaAd, int64, error) {
	var adsList []MetaAd
	var total int64

	query := r.db.Model(&MetaAd{})

	if filter.CampaignID != "" {
		query = query.Where("campaign_id = ?", filter.CampaignID)
	}
	if filter.AdSetID != "" {
		query = query.Where("adset_id = ?", filter.AdSetID)
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

	err := query.Order("created_time DESC").Limit(filter.Limit).Offset(offset).Find(&adsList).Error
	return adsList, total, err
}

func (r *repository) FindCreativeRawPayload(creativeID string) (string, error) {
	var result struct {
		RawPayload string
	}
	err := r.db.Table("ad_creatives").Select("raw_payload").Where("creative_id = ? AND deleted_at IS NULL", creativeID).First(&result).Error
	return result.RawPayload, err
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
	UpdatedTime       *time.Time      `gorm:"column:updated_time"`
}
func (r *repository) FindAdDashboard(filter AdFilter) ([]adDashboardScan, int64, error) {
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
  a.creative_id,
  a.updated_time
FROM meta_ads a
JOIN meta_campaigns c ON a.campaign_id = c.id
JOIN meta_ad_sets s ON a.adset_id = s.id
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

func javaStrToInt(s string) (int, error) {
	var v int
	_, _ = fmt.Sscanf(s, "%d", &v)
	return v, nil
}

func intToStr(v int) string {
	return fmt.Sprintf("%d", v)
}
