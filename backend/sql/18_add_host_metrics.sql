-- 18_add_host_metrics.sql
-- 添加机器监控指标表
-- @author Claude
-- @description 存储机器的CPU、内存、磁盘等监控数据
-- @modified 2026-02-06

CREATE TABLE IF NOT EXISTS host_metrics (
    id BIGSERIAL PRIMARY KEY,
    host_id VARCHAR(255) NOT NULL,

    -- CPU 指标
    cpu_usage_percent DECIMAL(5,2),
    cpu_cores_used DECIMAL(10,2),

    -- 内存指标
    memory_total_gb BIGINT,
    memory_used_gb BIGINT,
    memory_usage_percent DECIMAL(5,2),

    -- 磁盘指标
    disk_total_gb BIGINT,
    disk_used_gb BIGINT,
    disk_usage_percent DECIMAL(5,2),

    -- 网络指标（可选）
    network_rx_bytes BIGINT,
    network_tx_bytes BIGINT,

    -- 采集时间
    collected_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- 索引
    CONSTRAINT fk_host_metrics_host FOREIGN KEY (host_id) REFERENCES hosts(id) ON DELETE CASCADE
);

-- 创建索引以提高查询性能
CREATE INDEX idx_host_metrics_host_id ON host_metrics(host_id);
CREATE INDEX idx_host_metrics_collected_at ON host_metrics(collected_at);
CREATE INDEX idx_host_metrics_host_time ON host_metrics(host_id, collected_at DESC);

-- 添加表注释
COMMENT ON TABLE host_metrics IS '机器监控指标数据';
COMMENT ON COLUMN host_metrics.host_id IS '机器ID';
COMMENT ON COLUMN host_metrics.cpu_usage_percent IS 'CPU使用率(%)';
COMMENT ON COLUMN host_metrics.memory_usage_percent IS '内存使用率(%)';
COMMENT ON COLUMN host_metrics.disk_usage_percent IS '磁盘使用率(%)';
COMMENT ON COLUMN host_metrics.collected_at IS '采集时间';
