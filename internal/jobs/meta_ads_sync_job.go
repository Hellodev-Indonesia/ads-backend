package jobs

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/alex/ads_backend/internal/meta/ad_account"
	adcreative "github.com/alex/ads_backend/internal/meta/ad_creative"
	"github.com/alex/ads_backend/internal/meta/ads"
	"github.com/alex/ads_backend/internal/meta/adset"
	"github.com/alex/ads_backend/internal/meta/campaign"
	"github.com/alex/ads_backend/internal/meta/insight"
	insightDto "github.com/alex/ads_backend/internal/meta/insight/dto"
	metasync "github.com/alex/ads_backend/internal/meta/sync"
	syncDto "github.com/alex/ads_backend/internal/meta/sync/dto"
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
	Percentage uint8  `json:"percentage,omitempty"`
}

var stepLabels = map[string]string{
	metasync.SyncTypeAdAccounts:       "ad accounts",
	metasync.SyncTypeCampaigns:        "campaigns",
	metasync.SyncTypeAdsets:           "ad sets",
	metasync.SyncTypeAds:              "ads",
	metasync.SyncTypeAdCreatives:      "ad creatives",
	metasync.SyncTypeCampaignInsights: "campaign insights",
	metasync.SyncTypeAdInsights:       "ad insights",
	metasync.SyncTypeBusinesses:       "businesses",
}

type MetaAdsSyncJob struct {
	adAccountService  ad_account.Service
	campaignService   campaign.Service
	adSetService      adset.Service
	adsService        ads.Service
	insightService    insight.Service
	adCreativeService adcreative.Service
	syncLogService    *metasync.Service
	publisher         Publisher
	running           atomic.Bool
}

func NewMetaAdsSyncJob(
	adAccountService ad_account.Service,
	campaignService campaign.Service,
	adSetService adset.Service,
	adsService ads.Service,
	insightService insight.Service,
	syncLogService *metasync.Service,
	publisher Publisher,
	adCreativeService adcreative.Service,
) *MetaAdsSyncJob {
	return &MetaAdsSyncJob{
		adAccountService:  adAccountService,
		campaignService:   campaignService,
		adSetService:      adSetService,
		adsService:        adsService,
		insightService:    insightService,
		adCreativeService: adCreativeService,
		syncLogService:    syncLogService,
		publisher:         publisher,
	}
}

// Start creates sync batches and launches the job in the background.
// Returns metasync.ErrAlreadyRunning if a sync is currently in progress.
func (j *MetaAdsSyncJob) Start(ctx context.Context, req syncDto.TriggerSyncRequest) ([]*metasync.MetaSyncBatch, error) {
	if !j.running.CompareAndSwap(false, true) {
		return nil, metasync.ErrAlreadyRunning
	}

	var accountIDs []string
	if req.AdAccountID != "" {
		accountIDs = []string{req.AdAccountID}
	} else {
		// Fetch all active ad accounts
		// We use an empty filter to get all, then filter active. In production, we'd add an IsActive filter to the repo.
		accounts, _, err := j.adAccountService.GetAdAccounts(ad_account.AdAccountFilter{Limit: 1000})
		if err != nil {
			j.running.Store(false)
			return nil, fmt.Errorf("failed to fetch ad accounts: %v", err)
		}
		for _, acc := range accounts {
			if acc.IsActive {
				accountIDs = append(accountIDs, acc.ID)
			}
		}
	}

	// If len(accountIDs) is 0, we don't fail here. We will fetch them in execute after SyncAdAccounts.

	datePreset := "last_30d"
	var dp *string
	if req.DateStart == "" && req.DateStop == "" {
		dp = &datePreset
	}

	var parsedDateStart, parsedDateStop *time.Time
	if req.DateStart != "" {
		if t, err := time.Parse("2006-01-02", req.DateStart); err == nil {
			parsedDateStart = &t
		} else {
			log.Printf("Invalid DateStart format, ignoring: %v", err)
		}
	}
	if req.DateStop != "" {
		if t, err := time.Parse("2006-01-02", req.DateStop); err == nil {
			parsedDateStop = &t
		} else {
			log.Printf("Invalid DateStop format, ignoring: %v", err)
		}
	}

	syncScope := "incremental"
	if req.AdAccountID == "" {
		syncScope = "global"
	}

	batch, err := j.syncLogService.StartBatch(ctx, metasync.StartBatchInput{
		AdAccountID: req.AdAccountID,
		SyncMode:    "manual",
		SyncScope:   syncScope,
		DatePreset:  dp,
		DateStart:   parsedDateStart,
		DateStop:    parsedDateStop,
	})
	if err != nil {
		j.running.Store(false)
		return nil, fmt.Errorf("failed to start batch: %v", err)
	}

	j.publish(ctx, syncEvent{
		Event:     "sync:started",
		Message:   "Meta Ads sync started",
		BatchID:   batch.ID,
		BatchCode: batch.BatchCode,
	})

	batches := []*metasync.MetaSyncBatch{batch}

	insightReq := insightDto.SyncInsightRequest{
		DateStart: req.DateStart,
		DateStop:  req.DateStop,
	}

	go j.executeAll(batches, false, insightReq, accountIDs)

	return batches, nil
}

// IsRunning reports whether a sync is currently in progress.
func (j *MetaAdsSyncJob) IsRunning() bool {
	return j.running.Load()
}

func (j *MetaAdsSyncJob) executeAll(batches []*metasync.MetaSyncBatch, insightsOnly bool, insightReq insightDto.SyncInsightRequest, accountIDs []string) {
	defer j.running.Store(false)

	for _, batch := range batches {
		j.execute(batch, insightsOnly, insightReq, accountIDs)
	}
}

func (j *MetaAdsSyncJob) execute(batch *metasync.MetaSyncBatch, insightsOnly bool, insightReq insightDto.SyncInsightRequest, accountIDs []string) {

	ctx := context.Background()
	startTime := time.Now()

	hasError := false
	var firstError error
	var adAccountsCount, campaignCount, adSetCount, adsCount, adCreativeCount, campaignInsightCount, adInsightCount int

	totalSteps := 0
	if !insightsOnly {
		totalSteps += 1
	}
	stepsPerAccount := 0
	if !insightsOnly {
		stepsPerAccount += 4
	}
	if insightReq.Level == "" || insightReq.Level == "campaign" {
		stepsPerAccount += 1
	}
	if insightReq.Level == "" || insightReq.Level == "ad" || insightReq.Level == "adset" {
		stepsPerAccount += 1
	}
	totalSteps += len(accountIDs) * stepsPerAccount
	currentStep := 0

	if !insightsOnly {
		c, err := j.runSyncStep(
			ctx, batch.ID,
			metasync.SyncTypeAdAccounts,
			"/me/adaccounts",
			&currentStep, totalSteps,
			func() (int, error) { return j.adAccountService.SyncAdAccounts() },
		)
		adAccountsCount += c
		if err != nil {
			hasError = true
			firstError = setFirstError(firstError, err)
		}

		// Re-fetch active ad accounts if we didn't specify one
		if len(accountIDs) == 0 {
			accounts, _, _ := j.adAccountService.GetAdAccounts(ad_account.AdAccountFilter{Limit: 1000})
			for _, acc := range accounts {
				if acc.IsActive {
					accountIDs = append(accountIDs, acc.ID)
				}
			}
			totalSteps += len(accountIDs) * stepsPerAccount
		}
	}

	for idx, actID := range accountIDs {
		pct := uint8((float64(idx) / float64(len(accountIDs))) * 100)
		_ = j.syncLogService.UpdateBatchProgress(ctx, batch.ID, pct)
		j.publish(ctx, syncEvent{
			Event:      "sync:progress",
			Message:    fmt.Sprintf("Syncing ad account %d of %d", idx+1, len(accountIDs)),
			BatchID:    batch.ID,
			Percentage: pct,
		})

		accountHasError := false

		if !insightsOnly {
			c, err := j.runSyncStep(
				ctx, batch.ID,
				metasync.SyncTypeCampaigns,
				fmt.Sprintf("/%s/campaigns", actID),
				&currentStep, totalSteps,
				func() (int, error) { return j.campaignService.SyncCampaigns(actID) },
			)
			campaignCount += c
			if err != nil {
				accountHasError = true
				hasError = true
				firstError = setFirstError(firstError, err)
			}
		}

		if !accountHasError && !insightsOnly {
			c, err := j.runSyncStep(
				ctx, batch.ID,
				metasync.SyncTypeAdsets,
				fmt.Sprintf("/%s/adsets", actID),
				&currentStep, totalSteps,
				func() (int, error) { return j.adSetService.SyncAdSets(actID) },
			)
			adSetCount += c
			if err != nil {
				accountHasError = true
				hasError = true
				firstError = setFirstError(firstError, err)
			}
		}

		var syncedAds []ads.MetaAd
		if !accountHasError && !insightsOnly {
			c, err := j.runSyncStep(
				ctx, batch.ID,
				metasync.SyncTypeAds,
				fmt.Sprintf("/%s/ads", actID),
				&currentStep, totalSteps,
				func() (int, error) {
					count, models, e := j.adsService.SyncAdsWithList(actID)
					if e == nil {
						syncedAds = models
					}
					return count, e
				},
			)
			adsCount += c
			if err != nil {
				accountHasError = true
				hasError = true
				firstError = setFirstError(firstError, err)
			}
		}

		if !accountHasError && !insightsOnly {
			adRecords := make([]adcreative.AdRecord, 0, len(syncedAds))
			for _, a := range syncedAds {
				if a.CreativeID != "" {
					adRecords = append(adRecords, adcreative.AdRecord{
						ID:         a.ID,
						CreativeID: a.CreativeID,
						AdSetID:    a.AdSetID,
						CampaignID: a.CampaignID,
					})
				}
			}

			c, err := j.runSyncStep(
				ctx, batch.ID,
				metasync.SyncTypeAdCreatives,
				fmt.Sprintf("/%s/creatives", actID),
				&currentStep, totalSteps,
				func() (int, error) { return j.adCreativeService.SyncCreatives(actID, adRecords) },
			)
			adCreativeCount += c
			if err != nil {
				accountHasError = true
				hasError = true
				firstError = setFirstError(firstError, err)
			}
		}

		req := insightReq
		req.AdAccountID = actID
		if !insightsOnly {
			if req.DateStart == "" && req.DateStop == "" && req.DatePreset == "" {
				req.DatePreset = "last_30d"
			}
			req.TimeIncrement = 1
		}

		if (req.Level == "" || req.Level == "campaign") && !accountHasError {
			c, err := j.runSyncStep(
				ctx, batch.ID,
				metasync.SyncTypeCampaignInsights,
				fmt.Sprintf("/%s/insights?level=campaign", actID),
				&currentStep, totalSteps,
				func() (int, error) { return j.insightService.SyncCampaignInsights(req) },
			)
			campaignInsightCount += c
			if err != nil {
				accountHasError = true
				hasError = true
				firstError = setFirstError(firstError, err)
			}
		}

		if (req.Level == "" || req.Level == "ad" || req.Level == "adset") && !accountHasError {
			c, err := j.runSyncStep(
				ctx, batch.ID,
				metasync.SyncTypeAdInsights,
				fmt.Sprintf("/%s/insights?level=%s", actID, req.Level),
				&currentStep, totalSteps,
				func() (int, error) { return j.insightService.SyncAdInsights(req) },
			)
			adInsightCount += c
			if err != nil {
				accountHasError = true
				hasError = true
				firstError = setFirstError(firstError, err)
			}
		}
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
		"Meta Ads sync finished in %s (ad_accounts: %d, campaigns: %d, adsets: %d, ads: %d, ad_creatives: %d, campaign_insights: %d, ad_insights: %d)",
		elapsed, adAccountsCount, campaignCount, adSetCount, adsCount, adCreativeCount, campaignInsightCount, adInsightCount,
	)
}

func (j *MetaAdsSyncJob) runSyncStep(
	ctx context.Context,
	batchID uint64,
	syncType string,
	endpoint string,
	currentStep *int,
	totalSteps int,
	syncFunc func() (int, error),
) (int, error) {
	step, err := j.syncLogService.StartStep(ctx, batchID, syncType, endpoint)
	if err != nil {
		log.Printf("Failed to start sync step %s: %v", syncType, err)
		return 0, err
	}

	pct := uint8((float64(*currentStep) / float64(totalSteps)) * 100)

	j.publish(ctx, syncEvent{
		Event:      "sync:step:started",
		Message:    fmt.Sprintf("Syncing %s...", stepLabels[syncType]),
		BatchID:    batchID,
		Step:       syncType,
		Percentage: pct,
	})

	stepStart := time.Now()
	count, err := syncFunc()
	durationMs := time.Since(stepStart).Milliseconds()

	*currentStep++
	completedPct := uint8((float64(*currentStep) / float64(totalSteps)) * 100)

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
			Percentage: completedPct,
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
		Percentage: completedPct,
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
