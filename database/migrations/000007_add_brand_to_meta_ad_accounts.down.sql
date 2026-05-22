ALTER TABLE meta_ad_accounts
    DROP FOREIGN KEY fk_meta_ad_accounts_brand,
    DROP INDEX idx_meta_ad_accounts_brand_id,
    DROP COLUMN brand_id,
    DROP COLUMN currency,
    DROP COLUMN timezone_name,
    DROP COLUMN business_id,
    DROP COLUMN is_active;
