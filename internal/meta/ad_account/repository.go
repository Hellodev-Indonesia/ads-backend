package ad_account

import (
	"strings"

	"github.com/alex/ads_backend/internal/meta/ad_account/dto"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AdAccountFilter struct {
	Search        string
	Page          int
	Limit         int
	BrandID       *uint64
	BrandIDs      []uint64
	BusinessID    *string
	AccountStatus *int
	DateStart     string
	DateStop      string
}

type Repository interface {
	UpsertBatch(accounts []MetaAdAccount) error
	FindAll(filter AdAccountFilter) ([]MetaAdAccount, int64, error)
	FindUnassigned(filter AdAccountFilter) ([]MetaAdAccount, int64, error)
	FindByID(id string) (*MetaAdAccount, error)
	GetUniqueBusinesses() ([]dto.BusinessOptionResponse, error)
	Update(account *MetaAdAccount) error
	UpdateBrandID(id string, brandID *uint64) error
	UpdateBrandIDBatch(ids []string, brandID *uint64) error
	UpdateBrandIDByBusiness(businessID string, brandID *uint64) error
	FindBrandDashboard(filter AdAccountFilter) ([]brandDashboardScan, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) UpsertBatch(accounts []MetaAdAccount) error {
	if len(accounts) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"account_status",
			"currency",
			"timezone_name",
			"business_id",
			"business_name",
			"is_active",
			"synced_at",
		}),
	}).CreateInBatches(accounts, 100).Error
}

func (r *repository) FindAll(filter AdAccountFilter) ([]MetaAdAccount, int64, error) {
	var accounts []MetaAdAccount
	var total int64

	query := r.db.Model(&MetaAdAccount{})

	if filter.Search != "" {
		query = query.Where("name LIKE ?", "%"+filter.Search+"%")
	}
	
	if filter.BusinessID != nil {
		query = query.Where("business_id = ?", *filter.BusinessID)
	}

	if filter.BrandID != nil {
		query = query.Where("brand_id = ?", *filter.BrandID)
	}

	if filter.AccountStatus != nil {
		query = query.Where("account_status = ?", *filter.AccountStatus)
	}

	query.Count(&total)

	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit

	err := query.Order("name ASC").Limit(filter.Limit).Offset(offset).Find(&accounts).Error
	return accounts, total, err
}

func (r *repository) FindByID(id string) (*MetaAdAccount, error) {
	var account MetaAdAccount
	err := r.db.First(&account, "id = ?", id).Error
	return &account, err
}

func (r *repository) Update(account *MetaAdAccount) error {
	return r.db.Save(account).Error
}

func (r *repository) UpdateBrandID(id string, brandID *uint64) error {
	return r.db.Model(&MetaAdAccount{}).Where("id = ?", id).Update("brand_id", brandID).Error
}

func (r *repository) UpdateBrandIDBatch(ids []string, brandID *uint64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Model(&MetaAdAccount{}).Where("id IN ?", ids).Update("brand_id", brandID).Error
}

func (r *repository) GetUniqueBusinesses() ([]dto.BusinessOptionResponse, error) {
	var businesses []dto.BusinessOptionResponse
	err := r.db.Model(&MetaAdAccount{}).
		Select("business_id, business_name").
		Where("business_id IS NOT NULL AND business_id != ''").
		Group("business_id, business_name").
		Find(&businesses).Error
	return businesses, err
}

func (r *repository) UpdateBrandIDByBusiness(businessID string, brandID *uint64) error {
	return r.db.Model(&MetaAdAccount{}).Where("business_id = ?", businessID).Update("brand_id", brandID).Error
}

func (r *repository) FindUnassigned(filter AdAccountFilter) ([]MetaAdAccount, int64, error) {
	var accounts []MetaAdAccount
	var total int64

	query := r.db.Model(&MetaAdAccount{}).Where("brand_id IS NULL")

	if filter.Search != "" {
		query = query.Where("name LIKE ?", "%"+filter.Search+"%")
	}

	if filter.BusinessID != nil {
		query = query.Where("business_id = ?", *filter.BusinessID)
	}

	if filter.AccountStatus != nil {
		query = query.Where("account_status = ?", *filter.AccountStatus)
	}

	query.Count(&total)

	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit

	err := query.Order("name ASC").Limit(filter.Limit).Offset(offset).Find(&accounts).Error
	return accounts, total, err
}

type brandDashboardScan struct {
	BrandID             uint64   `gorm:"column:brand_id"`
	BrandName           string   `gorm:"column:brand_name"`
	BrandSlug           string   `gorm:"column:brand_slug"`
	BrandPhoto          *string  `gorm:"column:brand_photo"`
	AdAccountCount      int      `gorm:"column:ad_account_count"`
	ActiveCampaignCount int      `gorm:"column:active_campaign_count"`
	TotalSpends         *float64 `gorm:"column:total_spends"`
}
func (r *repository) FindBrandDashboard(filter AdAccountFilter) ([]brandDashboardScan, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = 25
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	where := "WHERE b.deleted_at IS NULL"
	var args []interface{}

	if filter.Search != "" {
		where += " AND b.name LIKE ?"
		args = append(args, "%"+filter.Search+"%")
	}

	if len(filter.BrandIDs) > 0 {
		var placeHolders []string
		for _, id := range filter.BrandIDs {
			placeHolders = append(placeHolders, "?")
			args = append(args, id)
		}
		where += " AND b.id IN (" + strings.Join(placeHolders, ",") + ")"
	}

	var total int64
	countSQL := "SELECT COUNT(*) FROM brands b " + where
	if err := r.db.Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit

	insightWhere := ""
	var insightArgs []interface{}
	if filter.DateStart != "" {
		insightWhere += " AND i.date_start >= ?"
		insightArgs = append(insightArgs, filter.DateStart)
	}
	if filter.DateStop != "" {
		insightWhere += " AND i.date_stop <= ?"
		insightArgs = append(insightArgs, filter.DateStop)
	}

	sql := `
SELECT 
    b.id AS brand_id,
    b.name AS brand_name,
    b.slug AS brand_slug,
    b.photo AS brand_photo,
    COUNT(DISTINCT a.id) AS ad_account_count,
    COUNT(DISTINCT CASE WHEN c.effective_status = 'ACTIVE' THEN c.id END) AS active_campaign_count,
    SUM(i.spend) AS total_spends
FROM brands b
LEFT JOIN meta_ad_accounts a ON b.id = a.brand_id
LEFT JOIN meta_campaigns c ON a.id COLLATE utf8mb4_unicode_ci = c.account_id COLLATE utf8mb4_unicode_ci
LEFT JOIN meta_insights i ON c.id COLLATE utf8mb4_unicode_ci = i.campaign_id COLLATE utf8mb4_unicode_ci AND i.level = 'campaign'` + insightWhere + `
` + where + `
GROUP BY b.id, b.name, b.slug, b.photo
ORDER BY b.name ASC
LIMIT ? OFFSET ?
`

	var finalArgs []interface{}
	finalArgs = append(finalArgs, insightArgs...)
	finalArgs = append(finalArgs, args...)
	finalArgs = append(finalArgs, filter.Limit, offset)

	var rows []brandDashboardScan
	if err := r.db.Raw(sql, finalArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}
