ALTER TABLE fraud_logs ADD COLUMN resolved_by BIGINT UNSIGNED NULL;
ALTER TABLE fraud_logs ADD CONSTRAINT fk_fraud_logs_resolved_by FOREIGN KEY (resolved_by) REFERENCES users(id) ON DELETE SET NULL;
