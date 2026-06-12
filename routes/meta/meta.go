package meta

import (
	"github.com/alex/ads_backend/internal/meta/ad_account"
	"github.com/alex/ads_backend/internal/meta/ads"
	"github.com/alex/ads_backend/internal/meta/adset"
	"github.com/alex/ads_backend/internal/meta/campaign"
	"github.com/alex/ads_backend/internal/meta/insight"
	"github.com/alex/ads_backend/internal/meta/activity"
	"github.com/gin-gonic/gin"
)

func RegisterMetaRoutes(
	r *gin.RouterGroup,
	adAccountHandler *ad_account.Handler,
	campaignHandler *campaign.Handler,
	adSetHandler *adset.Handler,
	adsHandler *ads.Handler,
	insightHandler *insight.Handler,
	activityHandler *activity.Handler,
) {
	ad_account.RegisterRoutes(r, adAccountHandler)
	campaign.RegisterRoutes(r, campaignHandler)
	adset.RegisterRoutes(r, adSetHandler)
	ads.RegisterRoutes(r, adsHandler)
	insight.RegisterRoutes(r, insightHandler)
	activity.RegisterRoutes(r, activityHandler)
}
