package jobs

import (
	"context"
	"log"
	"time"

	"github.com/alex/ads_backend/internal/meta/sync"
	"github.com/alex/ads_backend/internal/meta/sync/dto"
)

type Scheduler struct {
	syncService *sync.Service
	syncJob     *MetaAdsSyncJob
}

func NewScheduler(syncService *sync.Service, syncJob *MetaAdsSyncJob) *Scheduler {
	return &Scheduler{
		syncService: syncService,
		syncJob:     syncJob,
	}
}

// Start runs the background scheduler ticker.
func (s *Scheduler) Start() {
	go func() {
		// Run every minute
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			s.runTick()
		}
	}()
}

func (s *Scheduler) runTick() {
	ctx := context.Background()

	// 1. Fetch configuration
	config, err := s.syncService.GetConfig(ctx)
	if err != nil {
		log.Printf("[Scheduler] Error fetching sync config: %v", err)
		return
	}

	// 2. Check if active
	if !config.IsActive || config.IntervalMinutes <= 0 {
		return
	}

	// 3. Check if job is currently running
	if s.syncJob.IsRunning() {
		return
	}

	// 4. Fetch the last completed sync
	lastSync, err := s.syncService.GetLastSyncBatch(ctx)
	if err != nil {
		// This can happen if no sync has ever run. We should trigger the first one.
		log.Printf("[Scheduler] No previous completed sync found. Triggering initial scheduled sync.")
		s.triggerScheduledSync(ctx)
		return
	}

	if lastSync.FinishedAt == nil {
		return
	}

	// 5. Compare time elapsed
	minutesElapsed := time.Since(*lastSync.FinishedAt).Minutes()
	if minutesElapsed >= float64(config.IntervalMinutes) {
		log.Printf("[Scheduler] Triggering scheduled sync: %f minutes elapsed (configured: %d)", minutesElapsed, config.IntervalMinutes)
		s.triggerScheduledSync(ctx)
	}
}

func (s *Scheduler) triggerScheduledSync(ctx context.Context) {
	// For scheduled sync, we pull data for TODAY based on requirements
	today := time.Now().Format("2006-01-02")
	
	req := dto.TriggerSyncRequest{
		DateStart: today,
		DateStop:  today,
	}
	
	_, err := s.syncJob.Start(ctx, req)
	if err != nil {
		log.Printf("[Scheduler] Failed to trigger sync: %v", err)
	}
}
