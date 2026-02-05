-- ============================================
-- RemoteGPU 用户机器添加与主机凭据补充
-- ============================================
-- 文件: 14_machine_enrollments.sql
-- 说明: 新增用户机器添加任务表，补充主机 SSH 凭据字段
-- 执行顺序: 14
-- ============================================

-- CodeX 2026-02-04: 若单独执行本脚本，确保触发器函数存在。
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- CodeX 2026-02-04: 为主机补充区域与 SSH 凭据字段
ALTER TABLE hosts
    ADD COLUMN IF NOT EXISTS region VARCHAR(64) DEFAULT 'default',
    ADD COLUMN IF NOT EXISTS ssh_username VARCHAR(128),
    ADD COLUMN IF NOT EXISTS ssh_password TEXT,
    ADD COLUMN IF NOT EXISTS ssh_key TEXT,
    ADD COLUMN IF NOT EXISTS needs_collect BOOLEAN DEFAULT FALSE;

-- 用户机器添加任务表
CREATE TABLE IF NOT EXISTS machine_enrollments (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    name VARCHAR(128),
    hostname VARCHAR(256),
    region VARCHAR(64) DEFAULT 'default',
    address VARCHAR(256) NOT NULL,
    ssh_port INT DEFAULT 22,
    ssh_username VARCHAR(128),
    ssh_password TEXT,
    ssh_key TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    error_message TEXT,
    host_id VARCHAR(64),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_machine_enrollments_customer ON machine_enrollments(customer_id);
CREATE INDEX IF NOT EXISTS idx_machine_enrollments_status ON machine_enrollments(status);

CREATE TRIGGER update_machine_enrollments_updated_at
    BEFORE UPDATE ON machine_enrollments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE machine_enrollments IS '用户机器添加任务表';
COMMENT ON COLUMN machine_enrollments.status IS '状态: pending, success, failed';
