package sync

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type StartBatchInput struct {
	AdAccountID   string
	AdAccountName *string
	SyncMode      string
	SyncScope     string
	DatePreset    *string
	DateStart     *time.Time
	DateStop      *time.Time
}

func (s *Service) StartBatch(ctx context.Context, input StartBatchInput) (*MetaSyncBatch, error) {
	now := time.Now()

	syncMode := input.SyncMode
	if syncMode == "" {
		syncMode = "scheduled"
	}
	syncScope := input.SyncScope
	if syncScope == "" {
		syncScope = "incremental"
	}

	batch := &MetaSyncBatch{
		BatchCode:     generateBatchCode(now),
		Provider:      "meta",
		AdAccountID:   input.AdAccountID,
		AdAccountName: input.AdAccountName,
		SyncMode:      syncMode,
		SyncScope:     syncScope,
		DatePreset:    input.DatePreset,
		DateStart:     input.DateStart,
		DateStop:      input.DateStop,
		Status:        StatusRunning,
		StartedAt:     &now,
	}

	if err := s.repo.CreateBatch(ctx, batch); err != nil {
		return nil, err
	}
	return batch, nil
}

func (s *Service) CompleteBatch(ctx context.Context, batchID uint64) error {
	batch, err := s.repo.FindBatchByID(ctx, batchID)
	if err != nil {
		return err
	}
	now := time.Now()
	batch.Status = StatusCompleted
	batch.FinishedAt = &now
	batch.DurationMs = calculateDurationMs(batch.StartedAt, batch.FinishedAt)
	batch.ErrorMessage = nil
	return s.repo.UpdateBatch(ctx, batch)
}

func (s *Service) MarkBatchPartialFailed(ctx context.Context, batchID uint64, errInput error) error {
	batch, err := s.repo.FindBatchByID(ctx, batchID)
	if err != nil {
		return err
	}
	msg := errInput.Error()
	batch.Status = StatusPartialFailed
	batch.ErrorMessage = &msg
	return s.repo.UpdateBatch(ctx, batch)
}

func (s *Service) StartStep(ctx context.Context, batchID uint64, syncType string, endpoint string) (*MetaSyncStep, error) {
	now := time.Now()
	step := &MetaSyncStep{
		BatchID:   batchID,
		SyncType:  syncType,
		Endpoint:  nullableString(endpoint),
		Status:    StatusRunning,
		StartedAt: &now,
	}
	if err := s.repo.CreateStep(ctx, step); err != nil {
		return nil, err
	}
	return step, nil
}

func (s *Service) CompleteStep(ctx context.Context, stepID uint64, counts StepCounts) error {
	step, err := s.repo.FindStepByID(ctx, stepID)
	if err != nil {
		return err
	}
	now := time.Now()
	step.Status = StatusCompleted
	step.FinishedAt = &now
	step.DurationMs = calculateDurationMs(step.StartedAt, step.FinishedAt)
	step.TotalRecords = counts.TotalRecords
	step.InsertedCount = counts.InsertedCount
	step.UpdatedCount = counts.UpdatedCount
	step.SkippedCount = counts.SkippedCount
	step.FailedCount = counts.FailedCount
	step.RequestCount = counts.RequestCount
	step.ErrorCode = nil
	step.ErrorMessage = nil
	return s.repo.UpdateStep(ctx, step)
}

func (s *Service) FailStep(ctx context.Context, stepID uint64, errInput error) error {
	step, err := s.repo.FindStepByID(ctx, stepID)
	if err != nil {
		return err
	}
	now := time.Now()
	msg := errInput.Error()
	step.Status = StatusFailed
	step.FinishedAt = &now
	step.DurationMs = calculateDurationMs(step.StartedAt, step.FinishedAt)
	step.ErrorMessage = &msg
	step.FailedCount++
	return s.repo.UpdateStep(ctx, step)
}

func (s *Service) ListBatches(ctx context.Context, page, limit int) ([]MetaSyncBatch, int64, error) {
	if limit <= 0 {
		limit = 25
	}
	if page <= 0 {
		page = 1
	}
	batches, err := s.repo.ListBatches(ctx, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.CountBatches(ctx)
	if err != nil {
		return nil, 0, err
	}
	return batches, total, nil
}

func (s *Service) GetBatchByID(ctx context.Context, id uint64) (*MetaSyncBatch, error) {
	return s.repo.FindBatchByID(ctx, id)
}

func (s *Service) RecalculateBatchSummary(ctx context.Context, batchID uint64) error {
	batch, err := s.repo.FindBatchByID(ctx, batchID)
	if err != nil {
		return err
	}
	counts, err := s.repo.SumStepsByBatchID(ctx, batchID)
	if err != nil {
		return err
	}
	failedSteps, err := s.repo.CountFailedStepsByBatchID(ctx, batchID)
	if err != nil {
		return err
	}
	batch.TotalRecords = counts.TotalRecords
	batch.InsertedCount = counts.InsertedCount
	batch.UpdatedCount = counts.UpdatedCount
	batch.SkippedCount = counts.SkippedCount
	batch.FailedCount = counts.FailedCount
	batch.RequestCount = counts.RequestCount
	if failedSteps > 0 {
		batch.Status = StatusPartialFailed
	}
	return s.repo.UpdateBatch(ctx, batch)
}

func generateBatchCode(t time.Time) string {
	return fmt.Sprintf("META-SYNC-%s-%d", t.Format("20060102-150405"), t.UnixNano())
}

func calculateDurationMs(startedAt *time.Time, finishedAt *time.Time) uint64 {
	if startedAt == nil || finishedAt == nil {
		return 0
	}
	return uint64(finishedAt.Sub(*startedAt).Milliseconds())
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func (s *Service) CleanupOrphanedBatches(ctx context.Context) error {
	if err := s.repo.FailOrphanedSteps(ctx); err != nil {
		return err
	}
	if err := s.repo.FailOrphanedBatches(ctx); err != nil {
		return err
	}
	return nil
}
