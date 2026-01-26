-- ============================================
-- RemoteGPU 监控数据表
-- ============================================
-- 文件: 07_monitoring.sql
-- 说明: 创建主机监控、GPU监控、环境监控相关表
-- 执行顺序: 7
-- ============================================

-- 主机监控数据表
CREATE TABLE IF NOT EXISTS host_metrics (
    id BIGSERIAL PRIMARY KEY,
    host_id VARCHAR(64) NOT NULL,
    cpu_usage_percent FLOAT,
    cpu_load_1m FLOAT,
    cpu_load_5m FLOAT,
    cpu_load_15m FLOAT,
    memory_used BIGINT,
    memory_available BIGINT,
    memory_usage_percent FLOAT,
    disk_used BIGINT,
    disk_available BIGINT,
    disk_usage_percent FLOAT,
    disk_io_read_bytes BIGINT,
    disk_io_write_bytes BIGINT,
    network_rx_bytes BIGINT,
    network_tx_bytes BIGINT,
    network_rx_packets BIGINT,
    network_tx_packets BIGINT,
    gpu_avg_utilization FLOAT,
    gpu_avg_memory_used BIGINT,
    gpu_avg_temperature FLOAT,
    gpu_avg_power FLOAT,
    collected_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_host_metrics_host_time ON host_metrics(host_id, collected_at DESC);
CREATE INDEX idx_host_metrics_collected_at ON host_metrics(collected_at DESC);

-- 添加注释
COMMENT ON TABLE host_metrics IS '主机监控数据表（时序数据）';
COMMENT ON COLUMN host_metrics.cpu_usage_percent IS 'CPU使用率(%)';
COMMENT ON COLUMN host_metrics.memory_usage_percent IS '内存使用率(%)';
COMMENT ON COLUMN host_metrics.disk_usage_percent IS '磁盘使用率(%)';

-- GPU监控数据表
CREATE TABLE IF NOT EXISTS gpu_metrics (
    id BIGSERIAL PRIMARY KEY,
    gpu_id BIGINT NOT NULL,
    host_id VARCHAR(64) NOT NULL,
    utilization_percent FLOAT,
    memory_used BIGINT,
    memory_usage_percent FLOAT,
    temperature FLOAT,
    power_draw FLOAT,
    fan_speed_percent FLOAT,
    sm_clock INT,
    memory_clock INT,
    process_count INT,
    collected_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_gpu_metrics_gpu_time ON gpu_metrics(gpu_id, collected_at DESC);
CREATE INDEX idx_gpu_metrics_host_time ON gpu_metrics(host_id, collected_at DESC);
CREATE INDEX idx_gpu_metrics_collected_at ON gpu_metrics(collected_at DESC);

-- 添加注释
COMMENT ON TABLE gpu_metrics IS 'GPU监控数据表（时序数据）';
COMMENT ON COLUMN gpu_metrics.utilization_percent IS 'GPU使用率(%)';
COMMENT ON COLUMN gpu_metrics.temperature IS '温度(℃)';
COMMENT ON COLUMN gpu_metrics.power_draw IS '功耗(瓦)';

-- 环境监控数据表
CREATE TABLE IF NOT EXISTS environment_metrics (
    id BIGSERIAL PRIMARY KEY,
    env_id VARCHAR(64) NOT NULL,
    cpu_usage_percent FLOAT,
    memory_usage_percent FLOAT,
    gpu_usage_percent FLOAT,
    disk_usage_percent FLOAT,
    network_rx_bytes BIGINT,
    network_tx_bytes BIGINT,
    collected_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_environment_metrics_env_time ON environment_metrics(env_id, collected_at DESC);
CREATE INDEX idx_environment_metrics_collected_at ON environment_metrics(collected_at DESC);

-- 添加注释
COMMENT ON TABLE environment_metrics IS '环境监控数据表（时序数据）';
