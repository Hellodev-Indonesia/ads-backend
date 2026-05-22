package jobs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/alex/ads_backend/config"
	"github.com/alex/ads_backend/internal/meta/ads"
	"github.com/alex/ads_backend/internal/meta/adset"
	"github.com/alex/ads_backend/internal/meta/campaign"
	"github.com/alex/ads_backend/internal/meta/insight"
	metasync "github.com/alex/ads_backend/internal/meta/sync"
)

// Publisher pushes real-time events to a channel (e.g. Centrifugo).
type Publisher interface {
	Publish(ctx context.Context, channel string, data any) error
}

type syncEvent struct {
	Event      string `json:"event"`
	Message    string `json:"message"`
	BatchID    uint64 `json:"batch_id,omitempty"`
	BatchCode  string `json:"batch_code,omitempty"`
	Step       string `json:"step,omitempty"`
	Count      int    `json:"count,omitempty"`
	DurationMs int64  `json:"duration_ms,omitempty"`
	Error      string `json:"error,omitempty"`
}

var stepLabels = map[string]string{
	metasync.SyncTypeCampaigns:       "campaigns",
	metasync.SyncTypeAdsets:          "ad sets",
	metasync.SyncTypeAds:             "ads",
	metasync.SyncTypeCampaignInsights: "campaign insights",
	metasync.SyncTypeAdInsights:      "ad insights",
}

type MetaAdsSyncJob struct {
	campaignService campaign.Service
	adSetService    adset.Service
	adsService      ads.Service
	insightService  insight.Service
	syncLogService  *metasync.Service
	publisher       Publisher
	running         atomic.Bool
}

func NewMetaAdsSyncJob(
	campaignService campaign.Service,
	adSetService adset.Service,
	adsService ads.Service,
	insightService insight.Service,
	syncLogService *metasync.Service,
	publisher Publisher,
) *MetaAdsSyncJob {
	return &MetaAdsSyncJob{
		campaignService: campaignService,
		adSetService:    adSetService,
		adsService:      adsService,
		insightService:  insightService,
		syncLogService:  syncLogService,
		publisher:       publisher,
	}
}

// Start creates a sync batch and launches the job in the background.
// Returns metasync.ErrAlreadyRunning if a sync is currently in progress.
func (j *MetaAdsSyncJob) Start(ctx context.Context) (*metasync.MetaSyncBatch, error) {
	if !j.running.CompareAndSwap(false, true) {
		return nil, metasync.ErrAlreadyRunning
	}

	adAccountID := config.MetaAdAccountID
	if adAccountID == "" {
		j.running.Store(false)
		return nil, errors.New("META_AD_ACCOUNT_ID is not configured")
	}

	datePreset := "last_30d"
	batch, err := j.syncLogService.StartBatch(ctx, metasync.StartBatchInput{
		AdAccountID: adAccountID,
		SyncMode:    "manual",
		SyncScope:   "incremental",
		DatePreset:  &datePreset,
	})
	if err != nil {
		j.running.Store(false)
		return nil, err
	}

	j.publish(ctx, syncEvent{
		Event:     "sync:started",
		Message:   "Meta Ads sync started",
		BatchID:   batch.ID,
		BatchCode: batch.BatchCode,
	})

	go j.execute(batch)

	return batch, nil
}

// IsRunning reports whether a sync is currently in progress.
func (j *MetaAdsSyncJob) IsRunning() bool {
	return j.running.Load()
}

func (j *MetaAdsSyncJob) execute(batch *metasync.MetaSyncBatch) {
	defer j.running.Store(false)

	ctx := context.Background()
	adAccountID := batch.AdAccountID
	startTime := time.Now()

	hasError := false
	var firstError error

	campaignCount, err := j.runSyncStep(
		ctx, batch.ID,
		metasync.SyncTypeCampaigns,
		fmt.Sprintf("/%s/campaigns", adAccountID),
		func() (int, error) { return j.campaignService.SyncCampaigns(adAccountID) },
	)
	if err != nil {
		hasError = true
		firstError = setFirstError(firstError, err)
	}

	adSetCount, err := j.runSyncStep(
		ctx, batch.ID,
		metasync.SyncTypeAdsets,
		fmt.Sprintf("/%s/adsets", adAccountID),
		func() (int, error) { return j.adSetService.SyncAdSets(adAccountID) },
	)
	if err != nil {
		hasError = true
		firstError = setFirstError(firstError, err)
	}

	adsCount, err := j.runSyncStep(
		ctx, batch.ID,
		metasync.SyncTypeAds,
		fmt.Sprintf("/%s/ads", adAccountID),
		func() (int, error) { return j.adsService.SyncAds(adAccountID) },
	)
	if err != nil {
		hasError = true
		firstError = setFirstError(firstError, err)
	}

	campaignInsightCount, err := j.runSyncStep(
		ctx, batch.ID,
		metasync.SyncTypeCampaignInsights,
		fmt.Sprintf("/%s/insights?level=campaign", adAccountID),
		func() (int, error) { return j.insightService.SyncCampaignInsights(adAccountID) },
	)
	if err != nil {
		hasError = true
		firstError = setFirstError(firstError, err)
	}

	adInsightCount, err := j.runSyncStep(
		ctx, batch.ID,
		metasync.SyncTypeAdInsights,
		fmt.Sprintf("/%s/insights?level=ad", adAccountID),
		func() (int, error) { return j.insightService.SyncAdInsights(adAccountID) },
	)
	if err != nil {
		hasError = true
		firstError = setFirstError(firstError, err)
	}

	if err := j.syncLogService.RecalculateBatchSummary(ctx, batch.ID); err != nil {
		log.Printf("Failed to recalculate batch summary: %v", err)
	}

	elapsed := time.Since(startTime)

	if hasError {
		if err := j.syncLogService.MarkBatchPartialFailed(ctx, batch.ID, firstError); err != nil {
			log.Printf("Failed to mark batch as partial failed: %v", err)
		}
		j.publish(ctx, syncEvent{
			Event:      "sync:partial_failed",
			Message:    fmt.Sprintf("Sync completed with errors in %.1fs", elapsed.Seconds()),
			BatchID:    batch.ID,
			BatchCode:  batch.BatchCode,
			DurationMs: elapsed.Milliseconds(),
			Error:      firstError.Error(),
		})
	} else {
		if err := j.syncLogService.CompleteBatch(ctx, batch.ID); err != nil {
			log.Printf("Failed to complete batch: %v", err)
		}
		j.publish(ctx, syncEvent{
			Event:      "sync:completed",
			Message:    fmt.Sprintf("All data synced successfully in %.1fs", elapsed.Seconds()),
			BatchID:    batch.ID,
			BatchCode:  batch.BatchCode,
			DurationMs: elapsed.Milliseconds(),
		})
	}

	log.Printf(
		"Meta Ads sync finished in %s (campaigns: %d, adsets: %d, ads: %d, campaign_insights: %d, ad_insights: %d)",
		elapsed, campaignCount, adSetCount, adsCount, campaignInsightCount, adInsightCount,
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

	j.publish(ctx, syncEvent{
		Event:   "sync:step:started",
		Message: fmt.Sprintf("Syncing %s...", stepLabels[syncType]),
		BatchID: batchID,
		Step:    syncType,
	})

	stepStart := time.Now()
	count, err := syncFunc()
	durationMs := time.Since(stepStart).Milliseconds()

	if err != nil {
		log.Printf("Error syncing %s: %v", syncType, err)
		if failErr := j.syncLogService.FailStep(ctx, step.ID, err); failErr != nil {
			log.Printf("Failed to mark sync step %s as failed: %v", syncType, failErr)
		}
		j.publish(ctx, syncEvent{
			Event:      "sync:step:failed",
			Message:    fmt.Sprintf("Failed to sync %s", stepLabels[syncType]),
			BatchID:    batchID,
			Step:       syncType,
			DurationMs: durationMs,
			Error:      err.Error(),
		})
		return count, err
	}

	if err := j.syncLogService.CompleteStep(ctx, step.ID, metasync.StepCounts{
		TotalRecords: toUint(count),
		RequestCount: 1,
	}); err != nil {
		log.Printf("Failed to complete sync step %s: %v", syncType, err)
		return count, err
	}

	j.publish(ctx, syncEvent{
		Event:      "sync:step:completed",
		Message:    fmt.Sprintf("Synced %d %s in %.1fs", count, stepLabels[syncType], float64(durationMs)/1000),
		BatchID:    batchID,
		Step:       syncType,
		Count:      count,
		DurationMs: durationMs,
	})

	log.Printf("Synced %d records for %s", count, syncType)
	return count, nil
}

func (j *MetaAdsSyncJob) publish(ctx context.Context, event syncEvent) {
	if j.publisher == nil {
		return
	}
	if err := j.publisher.Publish(ctx, metasync.Channel, event); err != nil {
		log.Printf("Failed to publish sync event %s: %v", event.Event, err)
	}
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
