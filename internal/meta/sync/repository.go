package sync

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateBatch(ctx context.Context, batch *MetaSyncBatch) error {
	return r.db.WithContext(ctx).Create(batch).Error
}

func (r *Repository) UpdateBatch(ctx context.Context, batch *MetaSyncBatch) error {
	return r.db.WithContext(ctx).Save(batch).Error
}

func (r *Repository) UpdateBatchProgress(ctx context.Context, id uint64, progress uint8) error {
	return r.db.WithContext(ctx).Model(&MetaSyncBatch{}).Where("id = ?", id).Update("progress_percent", progress).Error
}

func (r *Repository) FindBatchByID(ctx context.Context, id uint64) (*MetaSyncBatch, error) {
	var batch MetaSyncBatch
	err := r.db.WithContext(ctx).Preload("Steps").First(&batch, id).Error
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *Repository) ListBatches(ctx context.Context, limit int, offset int) ([]MetaSyncBatch, error) {
	var batches []MetaSyncBatch
	err := r.db.WithContext(ctx).
		Order("started_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&batches).Error
	return batches, err
}

func (r *Repository) GetLastSyncBatch(ctx context.Context) (*MetaSyncBatch, error) {
	var batch MetaSyncBatch
	err := r.db.WithContext(ctx).
		Where("status = ?", StatusCompleted).
		Order("finished_at DESC").
		First(&batch).Error
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *Repository) CreateStep(ctx context.Context, step *MetaSyncStep) error {
	return r.db.WithContext(ctx).Create(step).Error
}

func (r *Repository) UpdateStep(ctx context.Context, step *MetaSyncStep) error {
	return r.db.WithContext(ctx).Save(step).Error
}

func (r *Repository) FindStepByID(ctx context.Context, id uint64) (*MetaSyncStep, error) {
	var step MetaSyncStep
	err := r.db.WithContext(ctx).First(&step, id).Error
	if err != nil {
		return nil, err
	}
	return &step, nil
}

func (r *Repository) ListStepsByBatchID(ctx context.Context, batchID uint64) ([]MetaSyncStep, error) {
	var steps []MetaSyncStep
	err := r.db.WithContext(ctx).
		Where("batch_id = ?", batchID).
		Order("id ASC").
		Find(&steps).Error
	return steps, err
}

func (r *Repository) SumStepsByBatchID(ctx context.Context, batchID uint64) (*StepCounts, error) {
	var counts StepCounts
	err := r.db.WithContext(ctx).
		Model(&MetaSyncStep{}).
		Select(`
			COALESCE(SUM(total_records), 0) AS total_records,
			COALESCE(SUM(inserted_count), 0) AS inserted_count,
			COALESCE(SUM(updated_count), 0) AS updated_count,
			COALESCE(SUM(skipped_count), 0) AS skipped_count,
			COALESCE(SUM(failed_count), 0) AS failed_count,
			COALESCE(SUM(request_count), 0) AS request_count
		`).
		Where("batch_id = ?", batchID).
		Scan(&counts).Error
	if err != nil {
		return nil, err
	}
	return &counts, nil
}

func (r *Repository) CountBatches(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&MetaSyncBatch{}).Count(&count).Error
	return count, err
}

func (r *Repository) CountFailedStepsByBatchID(ctx context.Context, batchID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&MetaSyncStep{}).
		Where("batch_id = ? AND status = ?", batchID, StatusFailed).
		Count(&count).Error
	return count, err
}

func (r *Repository) FailOrphanedBatches(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Model(&MetaSyncBatch{}).
		Where("status = ?", StatusRunning).
		Updates(map[string]interface{}{
			"status":        StatusFailed,
			"error_message": "Interrupted by server restart",
		}).Error
}

func (r *Repository) FailOrphanedSteps(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Model(&MetaSyncStep{}).
		Where("status = ?", StatusRunning).
		Updates(map[string]interface{}{
			"status":        StatusFailed,
			"error_message": "Interrupted by server restart",
		}).Error
}
