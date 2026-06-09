DROP INDEX idx_brands_slug ON brands;
DROP INDEX idx_brands_name ON brands;
ALTER TABLE brands DROP COLUMN slug;
