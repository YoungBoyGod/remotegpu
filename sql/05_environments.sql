-- ============================================
-- RemoteGPU 环境管理表
-- ============================================
-- 文件: 05_environments.sql
-- 说明: 创建开发环境、端口映射相关表
-- 执行顺序: 5
-- ============================================

-- 开发环境表
CREATE TABLE IF NOT EXISTS environments (
    id VARCHAR(64) PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    workspace_id BIGINT,
    host_id VARCHAR(64) NOT NULL,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    image VARCHAR(256) NOT NULL,
    status VARCHAR(20) DEFAULT 'creating',

    -- 资源配置
    cpu INT NOT NULL,
    memory BIGINT NOT NULL,
    gpu INT DEFAULT 0,
    gpu_memory BIGINT,
    storage BIGINT,
    temp_storage BIGINT,

    -- 访问配置
    ssh_port INT,
    ssh_enabled BOOLEAN DEFAULT true,
    rdp_port INT,
    rdp_enabled BOOLEAN DEFAULT false,
    jupyter_port INT,
    jupyter_token VARCHAR(128),
    jupyter_enabled BOOLEAN DEFAULT true,
    tensorboard_port INT,
    tensorboard_enabled BOOLEAN DEFAULT false,
    web_terminal_enabled BOOLEAN DEFAULT true,

    -- 挂载配置 (JSONB格式: [{"id": 1, "path": "/gemini/data-1", "readonly": true}])
    mounted_datasets JSONB,
    mounted_models JSONB,

    -- 环境配置
    env_vars JSONB,
    config JSONB,

    -- 容器信息
    container_id VARCHAR(128),
    pod_name VARCHAR(128),

    -- 时间戳
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    started_at TIMESTAMP,
    stopped_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_environments_customer ON environments(customer_id);
CREATE INDEX idx_environments_workspace ON environments(workspace_id);
CREATE INDEX idx_environments_host ON environments(host_id);
CREATE INDEX idx_environments_status ON environments(status);
CREATE INDEX idx_environments_created_at ON environments(created_at DESC);
CREATE INDEX idx_environments_deleted_at ON environments(deleted_at);

-- 为JSONB字段创建GIN索引
CREATE INDEX idx_environments_mounted_datasets ON environments USING GIN (mounted_datasets);
CREATE INDEX idx_environments_mounted_models ON environments USING GIN (mounted_models);
CREATE INDEX idx_environments_env_vars ON environments USING GIN (env_vars);
CREATE INDEX idx_environments_config ON environments USING GIN (config);

-- 创建更新时间触发器
CREATE TRIGGER update_environments_updated_at
    BEFORE UPDATE ON environments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE environments IS '开发环境表';
COMMENT ON COLUMN environments.status IS '状态: creating-创建中, running-运行中, stopped-已停止, error-错误, deleting-删除中';
COMMENT ON COLUMN environments.cpu IS 'CPU核心数';
COMMENT ON COLUMN environments.memory IS '内存(字节)';
COMMENT ON COLUMN environments.gpu IS 'GPU数量';
COMMENT ON COLUMN environments.gpu_memory IS 'GPU显存(字节)';
COMMENT ON COLUMN environments.storage IS '持久化存储空间(字节)';
COMMENT ON COLUMN environments.temp_storage IS '临时存储空间(字节)';
COMMENT ON COLUMN environments.jupyter_token IS 'JupyterLab访问令牌';
COMMENT ON COLUMN environments.mounted_datasets IS '挂载的数据集列表(JSONB): [{"id": 1, "path": "/gemini/data-1", "readonly": true}]';
COMMENT ON COLUMN environments.mounted_models IS '挂载的模型列表(JSONB): [{"id": 1, "path": "/gemini/pretrain1", "readonly": true}]';
COMMENT ON COLUMN environments.env_vars IS '环境变量(JSONB): {"KEY": "value"}';
COMMENT ON COLUMN environments.config IS '其他配置(JSONB)';
COMMENT ON COLUMN environments.deleted_at IS '软删除时间';

-- 端口映射表
CREATE TABLE IF NOT EXISTS port_mappings (
    id BIGSERIAL PRIMARY KEY,
    env_id VARCHAR(64) NOT NULL,
    service_type VARCHAR(32) NOT NULL,
    external_port INT NOT NULL UNIQUE,
    internal_port INT NOT NULL,
    status VARCHAR(20) DEFAULT 'active',

    -- 生命周期管理
    allocated_at TIMESTAMP DEFAULT NOW(),
    last_accessed_at TIMESTAMP,
    auto_release_hours INT DEFAULT 48,
    released_at TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_port_mappings_env ON port_mappings(env_id);
CREATE INDEX idx_port_mappings_external_port ON port_mappings(external_port);
CREATE INDEX idx_port_mappings_status ON port_mappings(status);

-- 添加注释
COMMENT ON TABLE port_mappings IS '端口映射表';
COMMENT ON COLUMN port_mappings.service_type IS '服务类型: ssh, rdp, jupyter, tensorboard, custom';
COMMENT ON COLUMN port_mappings.status IS '状态: active-活跃, released-已释放';
COMMENT ON COLUMN port_mappings.last_accessed_at IS '最后访问时间';
COMMENT ON COLUMN port_mappings.auto_release_hours IS '自动释放时间(小时)，默认48小时不活动后释放';
