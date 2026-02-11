ALTER TABLE port_mappings ADD COLUMN IF NOT EXISTS proxy_id VARCHAR(64);
ALTER TABLE port_mappings ADD COLUMN IF NOT EXISTS target_host VARCHAR(256);
ALTER TABLE port_mappings ADD COLUMN IF NOT EXISTS target_port INT;
ALTER TABLE port_mappings ADD COLUMN IF NOT EXISTS protocol VARCHAR(10) DEFAULT 'tcp';
CREATE INDEX IF NOT EXISTS idx_port_mappings_proxy_id ON port_mappings(proxy_id);
