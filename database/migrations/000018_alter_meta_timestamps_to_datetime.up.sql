ALTER TABLE meta_campaigns
MODIFY start_time DATETIME NULL,
MODIFY stop_time DATETIME NULL,
MODIFY created_time DATETIME NULL,
MODIFY updated_time DATETIME NULL;

ALTER TABLE meta_ad_sets
MODIFY start_time DATETIME NULL,
MODIFY end_time DATETIME NULL,
MODIFY created_time DATETIME NULL,
MODIFY updated_time DATETIME NULL;

ALTER TABLE meta_ads
MODIFY created_time DATETIME NULL,
MODIFY updated_time DATETIME NULL;

ALTER TABLE meta_businesses
MODIFY created_time DATETIME NULL;
