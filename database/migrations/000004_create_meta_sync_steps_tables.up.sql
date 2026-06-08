CREATE TABLE meta_sync_steps (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    batch_id BIGINT UNSIGNED NOT NULL,

    sync_type VARCHAR(100) NOT NULL,
    endpoint VARCHAR(255) NULL,

    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',

    started_at TIMESTAMP NULL,
    finished_at TIMESTAMP NULL,
    duration_ms BIGINT UNSIGNED DEFAULT 0,

    total_records INT UNSIGNED DEFAULT 0,
    inserted_count INT UNSIGNED DEFAULT 0,
    updated_count INT UNSIGNED DEFAULT 0,
    skipped_count INT UNSIGNED DEFAULT 0,
    failed_count INT UNSIGNED DEFAULT 0,

    request_count INT UNSIGNED DEFAULT 0,

    cursor_before TEXT NULL,
    cursor_after TEXT NULL,
    has_next TINYINT(1) DEFAULT 0,

    error_code VARCHAR(100) NULL,
    error_message TEXT NULL,

    metadata JSON NULL,

    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_meta_sync_steps_batch_id
        FOREIGN KEY (batch_id)
        REFERENCES meta_sync_batches(id)
        ON DELETE CASCADE,

    INDEX idx_meta_sync_steps_batch_id (batch_id),
    INDEX idx_meta_sync_steps_sync_type (sync_type),
    INDEX idx_meta_sync_steps_status (status)
);