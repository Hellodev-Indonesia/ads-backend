CREATE TABLE ad_creatives (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    creative_id VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NULL,
    title TEXT NULL,
    body TEXT NULL,
    image_url TEXT NULL,
    video_url TEXT NULL,
    destination_url TEXT NULL,
    normalized_domain VARCHAR(255) NULL,
    url_hash VARCHAR(64) NULL,
    raw_payload JSON NULL,
    synced_at TIMESTAMP NULL,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_ad_creatives_creative_id ON ad_creatives (creative_id);
CREATE INDEX idx_ad_creatives_normalized_domain ON ad_creatives (normalized_domain);
CREATE INDEX idx_ad_creatives_url_hash ON ad_creatives (url_hash);
