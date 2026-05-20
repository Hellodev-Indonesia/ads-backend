package jobs

import (
	"log"
	"time"

	"github.com/alex/ads_backend/config"
	"github.com/alex/ads_backend/internal/meta/ads"
	"github.com/alex/ads_backend/internal/meta/adset"
	"github.com/alex/ads_backend/internal/meta/campaign"
	"github.com/alex/ads_backend/internal/meta/insight"
)

type MetaAdsSyncJob struct {
	campaignService campaign.Service
	adSetService    adset.Service
	adsService      ads.Service
	insightService  insight.Service
}

func NewMetaAdsSyncJob(
	campaignService campaign.Service,
	adSetService adset.Service,
	adsService ads.Service,
	insightService insight.Service,
) *MetaAdsSyncJob {
	return &MetaAdsSyncJob{
		campaignService: campaignService,
		adSetService:    adSetService,
		adsService:      adsService,
		insightService:  insightService,
	}
}

// Start launches a ticker-based job that runs every 15 minutes
func (j *MetaAdsSyncJob) Start() {
	log.Println("Meta Ads Sync Job scheduler initialized (Interval: 15 minutes)")

	// Run immediately on start in background
	go j.Run()

	ticker := time.NewTicker(15 * time.Minute)
	go func() {
		for range ticker.C {
			j.Run()
		}
	}()
}

// Run performs the actual synchronization logic
func (j *MetaAdsSyncJob) Run() {
	log.Println("Starting Meta Ads synchronization job...")

	adAccountID := config.MetaAdAccountID
	if adAccountID == "" {
		log.Println("Warning: META_AD_ACCOUNT_ID is empty, skipping background sync")
		return
	}

	startTime := time.Now()

	// 1. Sync Campaigns
	campaignCount, err := j.campaignService.SyncCampaigns(adAccountID)
	if err != nil {
		log.Printf("Error syncing campaigns: %v", err)
	} else {
		log.Printf("Synced %d campaigns", campaignCount)
	}

	// 2. Sync Ad Sets
	adSetCount, err := j.adSetService.SyncAdSets(adAccountID)
	if err != nil {
		log.Printf("Error syncing adsets: %v", err)
	} else {
		log.Printf("Synced %d adsets", adSetCount)
	}

	// 3. Sync Ads
	adsCount, err := j.adsService.SyncAds(adAccountID)
	if err != nil {
		log.Printf("Error syncing ads: %v", err)
	} else {
		log.Printf("Synced %d ads", adsCount)
	}

	// 4. Sync Campaign Insights
	campaignInsightCount, err := j.insightService.SyncCampaignInsights(adAccountID)
	if err != nil {
		log.Printf("Error syncing campaign insights: %v", err)
	} else {
		log.Printf("Synced %d campaign insights", campaignInsightCount)
	}

	// 5. Sync Ad Insights
	adInsightCount, err := j.insightService.SyncAdInsights(adAccountID)
	if err != nil {
		log.Printf("Error syncing ad insights: %v", err)
	} else {
		log.Printf("Synced %d ad insights", adInsightCount)
	}

	elapsed := time.Since(startTime)
	log.Printf("Meta Ads sync completed in %s (campaigns: %d, adsets: %d, ads: %d, campaign_insights: %d, ad_insights: %d)",
		elapsed, campaignCount, adSetCount, adsCount, campaignInsightCount, adInsightCount)
}
