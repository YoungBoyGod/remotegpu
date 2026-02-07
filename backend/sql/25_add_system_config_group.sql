-- 为 system_configs 表添加配置分组字段
ALTER TABLE system_configs ADD COLUMN IF NOT EXISTS config_group VARCHAR(64) NOT NULL DEFAULT 'general';

-- 创建分组索引
CREATE INDEX IF NOT EXISTS idx_system_configs_group ON system_configs(config_group);

-- 更新现有配置的分组
UPDATE system_configs SET config_group = 'network' WHERE config_key LIKE '%port_range%';
UPDATE system_configs SET config_group = 'system' WHERE config_key IN ('system_name', 'system_version');
UPDATE system_configs SET config_group = 'environment' WHERE config_key LIKE '%environment%';

COMMENT ON COLUMN system_configs.config_group IS '配置分组: general, system, network, environment 等';
