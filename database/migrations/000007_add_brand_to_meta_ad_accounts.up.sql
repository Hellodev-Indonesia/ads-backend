ALTER TABLE meta_ad_accounts
    ADD COLUMN brand_id BIGINT UNSIGNED NULL,
    ADD COLUMN currency VARCHAR(10) NULL,
    ADD COLUMN timezone_name VARCHAR(50) NULL,
    ADD COLUMN business_id VARCHAR(50) NULL,
    ADD COLUMN is_active TINYINT(1) NOT NULL DEFAULT 1,
    ADD CONSTRAINT fk_meta_ad_accounts_brand FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE SET NULL;

CREATE INDEX idx_meta_ad_accounts_brand_id ON meta_ad_accounts(brand_id);
