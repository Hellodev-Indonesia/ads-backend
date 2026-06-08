CREATE TABLE fraud_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    brand_id BIGINT UNSIGNED NULL,
    ad_account_id VARCHAR(100) NULL,
    campaign_id VARCHAR(100) NULL,
    adset_id VARCHAR(100) NULL,
    ad_id VARCHAR(100) NULL,
    creative_id VARCHAR(100) NULL,
    event_type VARCHAR(100) NOT NULL,
    severity VARCHAR(50) NOT NULL,
    old_value TEXT NULL,
    new_value TEXT NULL,
    matched_rule_id BIGINT UNSIGNED NULL,
    message TEXT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'open',
    detected_at TIMESTAMP NULL,
    resolved_at TIMESTAMP NULL,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_fraud_logs_brand FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_fraud_logs_brand_id ON fraud_logs (brand_id);
CREATE INDEX idx_fraud_logs_ad_account_id ON fraud_logs (ad_account_id);
CREATE INDEX idx_fraud_logs_campaign_id ON fraud_logs (campaign_id);
CREATE INDEX idx_fraud_logs_creative_id ON fraud_logs (creative_id);
CREATE INDEX idx_fraud_logs_status ON fraud_logs (status);

CREATE TABLE alerts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    fraud_log_id BIGINT UNSIGNED NULL,
    brand_id BIGINT UNSIGNED NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    severity VARCHAR(50) NOT NULL,
    is_read TINYINT(1) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_alerts_fraud_log FOREIGN KEY (fraud_log_id) REFERENCES fraud_logs(id) ON DELETE SET NULL,
    CONSTRAINT fk_alerts_brand FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_alerts_brand_id ON alerts (brand_id);
CREATE INDEX idx_alerts_is_read ON alerts (is_read);
