CREATE TABLE meta_activities (
  id VARCHAR(255) PRIMARY KEY,
  ad_account_id VARCHAR(255) NOT NULL,
  actor_id VARCHAR(255) NULL,
  actor_name VARCHAR(255) NULL,
  object_id VARCHAR(255) NULL,
  object_name VARCHAR(255) NULL,
  object_type VARCHAR(100) NULL,
  event_type VARCHAR(100) NULL,
  event_time DATETIME NULL,
  extra_data JSON NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_meta_activities_ad_account_id (ad_account_id),
  INDEX idx_meta_activities_event_time (event_time),
  INDEX idx_meta_activities_object_id (object_id)
);
