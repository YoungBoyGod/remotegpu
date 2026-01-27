-- ============================================
-- RemoteGPU 通知和日志表
-- ============================================
-- 文件: 10_notifications_and_logs.sql
-- 说明: 创建通知、审计日志、操作日志相关表
-- 执行顺序: 10
-- ============================================

-- 通知表
CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    title VARCHAR(256) NOT NULL,
    content TEXT,
    type VARCHAR(32) NOT NULL,
    level VARCHAR(20) DEFAULT 'info',
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_notifications_customer ON notifications(customer_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);

-- 添加注释
COMMENT ON TABLE notifications IS '通知表';
COMMENT ON COLUMN notifications.type IS '类型: system-系统, environment-环境, billing-计费, alert-告警';
COMMENT ON COLUMN notifications.level IS '级别: info-信息, warning-警告, error-错误';

-- 审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT,
    username VARCHAR(128),
    action VARCHAR(128) NOT NULL,
    resource_type VARCHAR(64),
    resource_id VARCHAR(128),
    ip_address VARCHAR(64),
    user_agent TEXT,
    request_method VARCHAR(10),
    request_path VARCHAR(512),
    status_code INT,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_audit_logs_customer ON audit_logs(customer_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- 添加注释
COMMENT ON TABLE audit_logs IS '审计日志表';
COMMENT ON COLUMN audit_logs.action IS '操作: create, update, delete, login, logout等';
COMMENT ON COLUMN audit_logs.resource_type IS '资源类型: environment, dataset, model等';
