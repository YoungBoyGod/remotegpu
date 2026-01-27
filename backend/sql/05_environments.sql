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
    cpu INT NOT NULL,
    memory BIGINT NOT NULL,
    gpu INT DEFAULT 0,
    storage BIGINT,
    ssh_port INT,
    rdp_port INT,
    jupyter_port INT,
    container_id VARCHAR(128),
    pod_name VARCHAR(128),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    started_at TIMESTAMP,
    stopped_at TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_environments_customer ON environments(customer_id);
CREATE INDEX idx_environments_workspace ON environments(workspace_id);
CREATE INDEX idx_environments_host ON environments(host_id);
CREATE INDEX idx_environments_status ON environments(status);
CREATE INDEX idx_environments_created_at ON environments(created_at DESC);

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
COMMENT ON COLUMN environments.storage IS '存储空间(字节)';

-- 端口映射表
CREATE TABLE IF NOT EXISTS port_mappings (
    id BIGSERIAL PRIMARY KEY,
    env_id VARCHAR(64) NOT NULL,
    service_type VARCHAR(32) NOT NULL,
    external_port INT NOT NULL UNIQUE,
    internal_port INT NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    allocated_at TIMESTAMP DEFAULT NOW(),
    released_at TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_port_mappings_env ON port_mappings(env_id);
CREATE INDEX idx_port_mappings_external_port ON port_mappings(external_port);
CREATE INDEX idx_port_mappings_status ON port_mappings(status);

-- 添加注释
COMMENT ON TABLE port_mappings IS '端口映射表';
COMMENT ON COLUMN port_mappings.service_type IS '服务类型: ssh, rdp, jupyter, custom';
COMMENT ON COLUMN port_mappings.status IS '状态: active-活跃, released-已释放';
