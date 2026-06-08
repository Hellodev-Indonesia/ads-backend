CREATE TABLE meta_sync_batches (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    batch_code VARCHAR(100) NOT NULL UNIQUE,
    provider VARCHAR(50) NOT NULL DEFAULT 'meta',

    ad_account_id VARCHAR(100) NOT NULL,
    ad_account_name VARCHAR(255) NULL,

    sync_mode VARCHAR(50) NOT NULL DEFAULT 'scheduled',
    sync_scope VARCHAR(50) NOT NULL DEFAULT 'incremental',

    date_preset VARCHAR(50) NULL,
    date_start DATE NULL,
    date_stop DATE NULL,

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
    rate_limit_hit TINYINT(1) DEFAULT 0,

    error_message TEXT NULL,
    metadata JSON NULL,

    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_meta_sync_batches_ad_account_id (ad_account_id),
    INDEX idx_meta_sync_batches_status (status),
    INDEX idx_meta_sync_batches_started_at (started_at)
);