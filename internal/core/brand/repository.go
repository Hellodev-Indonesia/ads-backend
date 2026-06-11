package brand

import (
	"strings"

	"github.com/alex/ads_backend/internal/core/brand/dto"
	"gorm.io/gorm"
)

type Repository interface {
	Create(brand *Brand) error
	Update(brand *Brand) error
	DeleteBySlug(slug string) error
	FindBySlug(slug string) (*Brand, error)
	FindAll(filter dto.BrandFilter) ([]Brand, int64, error)
	FindBrandDashboard(filter dto.BrandDashboardFilter) ([]brandDashboardScan, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func FilterBrand(filter dto.BrandFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter.Name != "" {
			db = db.Where("name LIKE ?", "%"+filter.Name+"%")
		}
		if filter.IsActive != nil {
			db = db.Where("is_active = ?", *filter.IsActive)
		}
		return db
	}
}

func (r *repository) Create(brand *Brand) error {
	return r.db.Create(brand).Error
}

func (r *repository) Update(brand *Brand) error {
	return r.db.Save(brand).Error
}

func (r *repository) DeleteBySlug(slug string) error {
	return r.db.Where("slug = ?", slug).Delete(&Brand{}).Error
}

func (r *repository) FindBySlug(slug string) (*Brand, error) {
	var brand Brand
	err := r.db.Select("brands.*, (SELECT count(id) FROM meta_ad_accounts WHERE brand_id = brands.id) as ad_account_count").Where("slug = ?", slug).First(&brand).Error
	return &brand, err
}

func (r *repository) FindAll(filter dto.BrandFilter) ([]Brand, int64, error) {
	var brands []Brand
	var total int64

	q := r.db.Model(&Brand{}).Scopes(FilterBrand(filter))
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

	err := q.Select("brands.*, (SELECT count(id) FROM meta_ad_accounts WHERE brand_id = brands.id) as ad_account_count").Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&brands).Error
	return brands, total, err
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
func (r *repository) FindBrandDashboard(filter dto.BrandDashboardFilter) ([]brandDashboardScan, int64, error) {
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
