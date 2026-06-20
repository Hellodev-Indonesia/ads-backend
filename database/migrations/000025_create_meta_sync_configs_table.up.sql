CREATE TABLE meta_sync_configs (
    id INT PRIMARY KEY AUTO_INCREMENT,
    interval_minutes INT NOT NULL DEFAULT 180,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Insert default configuration
INSERT INTO meta_sync_configs (id, interval_minutes, is_active) VALUES (1, 180, TRUE);
