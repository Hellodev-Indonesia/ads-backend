ALTER TABLE brands ADD COLUMN slug VARCHAR(255) NOT NULL DEFAULT '';
UPDATE brands SET slug = LOWER(REPLACE(name, ' ', '-'));
CREATE UNIQUE INDEX idx_brands_slug ON brands(slug);
CREATE UNIQUE INDEX idx_brands_name ON brands(name);
