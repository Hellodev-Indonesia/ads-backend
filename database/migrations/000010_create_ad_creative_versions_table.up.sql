CREATE TABLE ad_creative_versions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    creative_id VARCHAR(100) NOT NULL,
    ad_id VARCHAR(100) NULL,
    adset_id VARCHAR(100) NULL,
    campaign_id VARCHAR(100) NULL,
    ad_account_id VARCHAR(100) NULL,
    brand_id BIGINT UNSIGNED NULL,
    name VARCHAR(255) NULL,
    title TEXT NULL,
    body TEXT NULL,
    image_url TEXT NULL,
    video_url TEXT NULL,
    destination_url TEXT NULL,
    normalized_domain VARCHAR(255) NULL,
    url_hash VARCHAR(64) NULL,
    raw_payload JSON NULL,
    changed_fields JSON NULL,
    change_type VARCHAR(100) NULL,
    synced_at TIMESTAMP NULL,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_ad_creative_versions_brand FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_ad_creative_versions_creative_id ON ad_creative_versions (creative_id);
CREATE INDEX idx_ad_creative_versions_ad_id ON ad_creative_versions (ad_id);
CREATE INDEX idx_ad_creative_versions_campaign_id ON ad_creative_versions (campaign_id);
CREATE INDEX idx_ad_creative_versions_ad_account_id ON ad_creative_versions (ad_account_id);
CREATE INDEX idx_ad_creative_versions_brand_id ON ad_creative_versions (brand_id);
CREATE INDEX idx_ad_creative_versions_url_hash ON ad_creative_versions (url_hash);
CREATE INDEX idx_ad_creative_versions_created_at ON ad_creative_versions (created_at);
