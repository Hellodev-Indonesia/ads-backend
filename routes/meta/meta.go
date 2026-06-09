package meta

import (
	"github.com/alex/ads_backend/internal/meta/ad_account"
	"github.com/alex/ads_backend/internal/meta/ads"
	"github.com/alex/ads_backend/internal/meta/adset"
	"github.com/alex/ads_backend/internal/meta/business"
	"github.com/alex/ads_backend/internal/meta/campaign"
	"github.com/alex/ads_backend/internal/meta/dashboard"
	"github.com/alex/ads_backend/internal/meta/insight"
	"github.com/gin-gonic/gin"
)

func RegisterMetaRoutes(
	r *gin.RouterGroup,
	adAccHandler *ad_account.Handler,
	campaignHandler *campaign.Handler,
	adSetHandler *adset.Handler,
	adsHandler *ads.Handler,
	insightHandler *insight.Handler,
	dashboardHandler *dashboard.Handler,
	businessHandler *business.Handler,
) {
	ad_account.RegisterRoutes(r, adAccHandler)
	campaign.RegisterRoutes(r, campaignHandler)
	adset.RegisterRoutes(r, adSetHandler)
	ads.RegisterRoutes(r, adsHandler)
	insight.RegisterRoutes(r, insightHandler)
	dashboard.RegisterRoutes(r, dashboardHandler)
	business.RegisterRoutes(r, businessHandler)
}
