ALTER TABLE fraud_logs 
ADD COLUMN actor_id VARCHAR(255) NULL AFTER event_type,
ADD COLUMN actor_name VARCHAR(255) NULL AFTER actor_id;
