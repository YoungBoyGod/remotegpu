-- 添加 Guacamole 支持字段

ALTER TABLE environments ADD COLUMN IF NOT EXISTS use_guacamole BOOLEAN DEFAULT false;
COMMENT ON COLUMN environments.use_guacamole IS '是否使用 Apache Guacamole 进行远程访问';

ALTER TABLE environments ADD COLUMN IF NOT EXISTS guacamole_conn_id VARCHAR(64);
COMMENT ON COLUMN environments.guacamole_conn_id IS 'Guacamole 连接 ID';

CREATE INDEX IF NOT EXISTS idx_environments_use_guacamole ON environments(use_guacamole);
