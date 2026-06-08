CREATE TABLE brand_whitelist_rules (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    brand_id BIGINT UNSIGNED NOT NULL,
    scope VARCHAR(50) NOT NULL,
    match_type VARCHAR(50) NOT NULL,
    value TEXT NOT NULL,
    normalized_value TEXT NULL,
    allow_subdomains TINYINT(1) NOT NULL DEFAULT 0,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    description TEXT NULL,
    created_by BIGINT UNSIGNED NULL,
    approved_by BIGINT UNSIGNED NULL,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_brand_whitelist_rules_brand FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_brand_whitelist_rules_brand_id ON brand_whitelist_rules (brand_id);
CREATE INDEX idx_brand_whitelist_rules_scope ON brand_whitelist_rules (scope);
CREATE INDEX idx_brand_whitelist_rules_match_type ON brand_whitelist_rules (match_type);
CREATE INDEX idx_brand_whitelist_rules_is_active ON brand_whitelist_rules (is_active);
