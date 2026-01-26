-- ============================================
-- RemoteGPU 系统配置表
-- ============================================
-- 文件: 02_system_config.sql
-- 说明: 创建系统配置相关表
-- 执行顺序: 2
-- ============================================

-- 系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
    id BIGSERIAL PRIMARY KEY,
    config_key VARCHAR(128) UNIQUE NOT NULL,
    config_value TEXT NOT NULL,
    config_type VARCHAR(32) NOT NULL DEFAULT 'string',
    description TEXT,
    is_public BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_system_configs_key ON system_configs(config_key);
CREATE INDEX idx_system_configs_public ON system_configs(is_public);

-- 创建更新时间触发器
CREATE TRIGGER update_system_configs_updated_at
    BEFORE UPDATE ON system_configs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE system_configs IS '系统配置表';
COMMENT ON COLUMN system_configs.config_key IS '配置键';
COMMENT ON COLUMN system_configs.config_value IS '配置值';
COMMENT ON COLUMN system_configs.config_type IS '值类型: string, integer, boolean, json';
COMMENT ON COLUMN system_configs.is_public IS '是否公开（前端可见）';

-- 插入默认配置
INSERT INTO system_configs (config_key, config_value, config_type, description, is_public) VALUES
('ssh_port_range_start', '30000', 'integer', 'SSH端口范围起始值', false),
('ssh_port_range_end', '40000', 'integer', 'SSH端口范围结束值', false),
('rdp_port_range_start', '40001', 'integer', 'RDP端口范围起始值', false),
('rdp_port_range_end', '50000', 'integer', 'RDP端口范围结束值', false),
('jupyter_port_range_start', '50001', 'integer', 'JupyterLab端口范围起始值', false),
('jupyter_port_range_end', '60000', 'integer', 'JupyterLab端口范围结束值', false),
('max_environments_per_user', '10', 'integer', '每个用户最大环境数', false),
('default_environment_timeout', '7200', 'integer', '默认环境超时时间（秒）', false),
('system_name', 'RemoteGPU', 'string', '系统名称', true),
('system_version', '1.0.0', 'string', '系统版本', true)
ON CONFLICT (config_key) DO NOTHING;
