package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alex/ads_backend/config"
	"github.com/alex/ads_backend/internal/meta/ads"
	"github.com/alex/ads_backend/internal/meta/adset"
	"github.com/alex/ads_backend/internal/meta/campaign"
	"github.com/alex/ads_backend/internal/meta/insight"
	"github.com/alex/ads_backend/internal/meta/sync_logs"
)

type MetaAdsSyncJob struct {
	campaignService campaign.Service
	adSetService    adset.Service
	adsService      ads.Service
	insightService  insight.Service
	syncLogService  *sync_logs.Service
}

func NewMetaAdsSyncJob(
	campaignService campaign.Service,
	adSetService adset.Service,
	adsService ads.Service,
	insightService insight.Service,
	syncLogService *sync_logs.Service,
) *MetaAdsSyncJob {
	return &MetaAdsSyncJob{
		campaignService: campaignService,
		adSetService:    adSetService,
		adsService:      adsService,
		insightService:  insightService,
		syncLogService:  syncLogService,
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

	ctx := context.Background()

	adAccountID := config.MetaAdAccountID
	if adAccountID == "" {
		log.Println("Warning: META_AD_ACCOUNT_ID is empty, skipping background sync")
		return
	}

	startTime := time.Now()

	datePreset := "last_30d"

	batch, err := j.syncLogService.StartBatch(ctx, sync_logs.StartBatchInput{
		AdAccountID: adAccountID,
		SyncMode:    "scheduled",
		SyncScope:   "incremental",
		DatePreset:  &datePreset,
	})

	if err != nil {
		log.Printf("Failed to create meta sync batch: %v", err)
		return
	}

	hasError := false
	var firstError error

	// 1. Sync Campaigns
	campaignCount, err := j.runSyncStep(
		ctx,
		batch.ID,
		sync_logs.SyncTypeCampaigns,
		fmt.Sprintf("/%s/campaigns", adAccountID),
		func() (int, error) {
			return j.campaignService.SyncCampaigns(adAccountID)
		},
	)

	if err != nil {
		hasError = true
		firstError = setFirstError(firstError, err)
	}

	// 2. Sync Ad Sets
	adSetCount, err := j.runSyncStep(
		ctx,
		batch.ID,
		sync_logs.SyncTypeAdsets,
		fmt.Sprintf("/%s/adsets", adAccountID),
		func() (int, error) {
			return j.adSetService.SyncAdSets(adAccountID)
		},
	)

	if err != nil {
		hasError = true
		firstError = setFirstError(firstError, err)
	}

	// 3. Sync Ads
	adsCount, err := j.runSyncStep(
		ctx,
		batch.ID,
		sync_logs.SyncTypeAds,
		fmt.Sprintf("/%s/ads", adAccountID),
		func() (int, error) {
			return j.adsService.SyncAds(adAccountID)
		},
	)

	if err != nil {
		hasError = true
		firstError = setFirstError(firstError, err)
	}

	// 4. Sync Campaign Insights
	campaignInsightCount, err := j.runSyncStep(
		ctx,
		batch.ID,
		sync_logs.SyncTypeCampaignInsights,
		fmt.Sprintf("/%s/insights?level=campaign", adAccountID),
		func() (int, error) {
			return j.insightService.SyncCampaignInsights(adAccountID)
		},
	)

	if err != nil {
		hasError = true
		firstError = setFirstError(firstError, err)
	}

	// 5. Sync Ad Insights
	adInsightCount, err := j.runSyncStep(
		ctx,
		batch.ID,
		sync_logs.SyncTypeAdInsights,
		fmt.Sprintf("/%s/insights?level=ad", adAccountID),
		func() (int, error) {
			return j.insightService.SyncAdInsights(adAccountID)
		},
	)

	if err != nil {
		hasError = true
		firstError = setFirstError(firstError, err)
	}

	if err := j.syncLogService.RecalculateBatchSummary(ctx, batch.ID); err != nil {
		log.Printf("Failed to recalculate batch summary: %v", err)
	}

	if hasError {
		if err := j.syncLogService.MarkBatchPartialFailed(ctx, batch.ID, firstError); err != nil {
			log.Printf("Failed to mark batch as partial failed: %v", err)
		}
	} else {
		if err := j.syncLogService.CompleteBatch(ctx, batch.ID); err != nil {
			log.Printf("Failed to complete batch: %v", err)
		}
	}

	elapsed := time.Since(startTime)

	log.Printf(
		"Meta Ads sync finished in %s (campaigns: %d, adsets: %d, ads: %d, campaign_insights: %d, ad_insights: %d)",
		elapsed,
		campaignCount,
		adSetCount,
		adsCount,
		campaignInsightCount,
		adInsightCount,
	)
}

func (j *MetaAdsSyncJob) runSyncStep(
	ctx context.Context,
	batchID uint64,
	syncType string,
	endpoint string,
	syncFunc func() (int, error),
) (int, error) {
	step, err := j.syncLogService.StartStep(ctx, batchID, syncType, endpoint)
	if err != nil {
		log.Printf("Failed to start sync step %s: %v", syncType, err)
		return 0, err
	}

	count, err := syncFunc()
	if err != nil {
		log.Printf("Error syncing %s: %v", syncType, err)

		if failErr := j.syncLogService.FailStep(ctx, step.ID, err); failErr != nil {
			log.Printf("Failed to mark sync step %s as failed: %v", syncType, failErr)
		}

		return count, err
	}

	err = j.syncLogService.CompleteStep(ctx, step.ID, sync_logs.StepCounts{
		TotalRecords: toUint(count),
		RequestCount: 1,
	})

	if err != nil {
		log.Printf("Failed to complete sync step %s: %v", syncType, err)
		return count, err
	}

	log.Printf("Synced %d records for %s", count, syncType)

	return count, nil
}

func setFirstError(current error, newErr error) error {
	if current != nil {
		return current
	}

	return newErr
}

func toUint(value int) uint {
	if value < 0 {
		return 0
	}

	return uint(value)
}
