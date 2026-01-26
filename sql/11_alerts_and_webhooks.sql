-- ============================================
-- RemoteGPU 告警和 Webhook 表
-- ============================================
-- 文件: 11_alerts_and_webhooks.sql
-- 说明: 创建告警规则、告警记录、Webhook 相关表
-- 执行顺序: 11
-- ============================================

-- 告警规则表
CREATE TABLE IF NOT EXISTS alert_rules (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    metric_type VARCHAR(64) NOT NULL,
    threshold FLOAT NOT NULL,
    comparison VARCHAR(10) NOT NULL,
    duration INT DEFAULT 60,
    severity VARCHAR(20) DEFAULT 'warning',
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_alert_rules_customer ON alert_rules(customer_id);
CREATE INDEX idx_alert_rules_enabled ON alert_rules(enabled);

-- 创建更新时间触发器
CREATE TRIGGER update_alert_rules_updated_at
    BEFORE UPDATE ON alert_rules
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE alert_rules IS '告警规则表';
COMMENT ON COLUMN alert_rules.metric_type IS '指标类型: cpu_usage, memory_usage, gpu_usage, disk_usage等';
COMMENT ON COLUMN alert_rules.comparison IS '比较运算符: >, <, >=, <=, ==';
COMMENT ON COLUMN alert_rules.duration IS '持续时间(秒)';
COMMENT ON COLUMN alert_rules.severity IS '严重程度: info, warning, critical';

-- 告警记录表
CREATE TABLE IF NOT EXISTS alert_records (
    id BIGSERIAL PRIMARY KEY,
    rule_id BIGINT NOT NULL,
    customer_id BIGINT,
    resource_type VARCHAR(64),
    resource_id VARCHAR(128),
    metric_type VARCHAR(64) NOT NULL,
    current_value FLOAT NOT NULL,
    threshold FLOAT NOT NULL,
    severity VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'firing',
    message TEXT,
    triggered_at TIMESTAMP DEFAULT NOW(),
    resolved_at TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_alert_records_rule ON alert_records(rule_id);
CREATE INDEX idx_alert_records_customer ON alert_records(customer_id);
CREATE INDEX idx_alert_records_status ON alert_records(status);
CREATE INDEX idx_alert_records_triggered_at ON alert_records(triggered_at DESC);

-- 添加注释
COMMENT ON TABLE alert_records IS '告警记录表';
COMMENT ON COLUMN alert_records.status IS '状态: firing-触发中, resolved-已解决';

-- Webhook 配置表
CREATE TABLE IF NOT EXISTS webhooks (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    name VARCHAR(256) NOT NULL,
    url VARCHAR(512) NOT NULL,
    secret VARCHAR(256),
    events TEXT[] NOT NULL,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_webhooks_customer ON webhooks(customer_id);
CREATE INDEX idx_webhooks_enabled ON webhooks(enabled);

-- 创建更新时间触发器
CREATE TRIGGER update_webhooks_updated_at
    BEFORE UPDATE ON webhooks
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE webhooks IS 'Webhook配置表';
COMMENT ON COLUMN webhooks.events IS '订阅的事件类型数组';

-- Webhook 调用日志表
CREATE TABLE IF NOT EXISTS webhook_logs (
    id BIGSERIAL PRIMARY KEY,
    webhook_id BIGINT NOT NULL,
    event_type VARCHAR(64) NOT NULL,
    payload JSONB,
    status_code INT,
    response_body TEXT,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_webhook_logs_webhook ON webhook_logs(webhook_id);
CREATE INDEX idx_webhook_logs_created_at ON webhook_logs(created_at DESC);

-- 添加注释
COMMENT ON TABLE webhook_logs IS 'Webhook调用日志表';
