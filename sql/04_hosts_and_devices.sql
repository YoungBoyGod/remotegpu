-- ============================================
-- RemoteGPU 主机和设备管理表
-- ============================================
-- 文件: 04_hosts_and_devices.sql
-- 说明: 创建主机、GPU设备相关表
-- 执行顺序: 4
-- ============================================

-- 主机表
CREATE TABLE IF NOT EXISTS hosts (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    hostname VARCHAR(256),
    ip_address VARCHAR(64) NOT NULL,
    public_ip VARCHAR(64),
    os_type VARCHAR(20) NOT NULL,
    os_version VARCHAR(64),
    arch VARCHAR(20) DEFAULT 'x86_64',
    deployment_mode VARCHAR(20) NOT NULL,
    k8s_node_name VARCHAR(128),
    status VARCHAR(20) DEFAULT 'offline',
    health_status VARCHAR(20) DEFAULT 'unknown',
    total_cpu INT NOT NULL,
    total_memory BIGINT NOT NULL,
    total_disk BIGINT,
    total_gpu INT DEFAULT 0,
    used_cpu INT DEFAULT 0,
    used_memory BIGINT DEFAULT 0,
    used_disk BIGINT DEFAULT 0,
    used_gpu INT DEFAULT 0,
    ssh_port INT DEFAULT 22,
    winrm_port INT,
    agent_port INT DEFAULT 8080,
    labels JSONB,
    tags TEXT[],
    last_heartbeat TIMESTAMP,
    registered_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_hosts_os_type ON hosts(os_type);
CREATE INDEX idx_hosts_status ON hosts(status);
CREATE INDEX idx_hosts_deployment_mode ON hosts(deployment_mode);
CREATE INDEX idx_hosts_labels ON hosts USING GIN(labels);
CREATE INDEX idx_hosts_health_status ON hosts(health_status);

-- 创建更新时间触发器
CREATE TRIGGER update_hosts_updated_at
    BEFORE UPDATE ON hosts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE hosts IS '主机表';
COMMENT ON COLUMN hosts.os_type IS '操作系统类型: linux, windows';
COMMENT ON COLUMN hosts.deployment_mode IS '部署模式: traditional-传统架构, kubernetes-K8s';
COMMENT ON COLUMN hosts.status IS '状态: online-在线, offline-离线, maintenance-维护中';
COMMENT ON COLUMN hosts.health_status IS '健康状态: healthy-健康, degraded-降级, unhealthy-不健康, unknown-未知';

-- GPU设备表
CREATE TABLE IF NOT EXISTS gpus (
    id BIGSERIAL PRIMARY KEY,
    host_id VARCHAR(64) NOT NULL,
    gpu_index INT NOT NULL,
    uuid VARCHAR(128) UNIQUE,
    name VARCHAR(128),
    brand VARCHAR(64),
    architecture VARCHAR(64),
    memory_total BIGINT,
    cuda_cores INT,
    compute_capability VARCHAR(32),
    status VARCHAR(20) DEFAULT 'available',
    health_status VARCHAR(20) DEFAULT 'healthy',
    allocated_to VARCHAR(64),
    allocated_at TIMESTAMP,
    power_limit INT,
    temperature_limit INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(host_id, gpu_index)
);

-- 创建索引
CREATE INDEX idx_gpus_host ON gpus(host_id);
CREATE INDEX idx_gpus_status ON gpus(status);
CREATE INDEX idx_gpus_allocated_to ON gpus(allocated_to);
CREATE INDEX idx_gpus_health_status ON gpus(health_status);

-- 创建更新时间触发器
CREATE TRIGGER update_gpus_updated_at
    BEFORE UPDATE ON gpus
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE gpus IS 'GPU设备表';
COMMENT ON COLUMN gpus.status IS '状态: available-可用, allocated-已分配, maintenance-维护中, error-错误';
COMMENT ON COLUMN gpus.health_status IS '健康状态: healthy-健康, degraded-降级, error-错误';
COMMENT ON COLUMN gpus.allocated_to IS '分配给的环境ID';
