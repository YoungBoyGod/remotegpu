-- ============================================
-- RemoteGPU 计费管理表
-- ============================================
-- 文件: 08_billing.sql
-- 说明: 创建计费记录、账单相关表
-- 执行顺序: 8
-- ============================================

-- 计费记录表
CREATE TABLE IF NOT EXISTS billing_records (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    env_id VARCHAR(64),
    resource_type VARCHAR(32) NOT NULL,
    quantity FLOAT NOT NULL,
    unit_price DECIMAL(10,4) NOT NULL,
    amount DECIMAL(10,4) NOT NULL,
    currency VARCHAR(10) DEFAULT 'CNY',
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_billing_records_customer ON billing_records(customer_id);
CREATE INDEX idx_billing_records_env ON billing_records(env_id);
CREATE INDEX idx_billing_records_start_time ON billing_records(start_time DESC);
CREATE INDEX idx_billing_records_created_at ON billing_records(created_at DESC);

-- 添加注释
COMMENT ON TABLE billing_records IS '计费记录表';
COMMENT ON COLUMN billing_records.resource_type IS '资源类型: cpu, memory, gpu, storage, network';
COMMENT ON COLUMN billing_records.quantity IS '数量（如CPU核心数、内存GB数）';
COMMENT ON COLUMN billing_records.unit_price IS '单价';
COMMENT ON COLUMN billing_records.amount IS '金额';

-- 账单表
CREATE TABLE IF NOT EXISTS invoices (
    id BIGSERIAL PRIMARY KEY,
    invoice_no VARCHAR(64) UNIQUE NOT NULL,
    customer_id BIGINT NOT NULL,
    billing_period_start TIMESTAMP NOT NULL,
    billing_period_end TIMESTAMP NOT NULL,
    total_amount DECIMAL(10,4) NOT NULL,
    currency VARCHAR(10) DEFAULT 'CNY',
    status VARCHAR(20) DEFAULT 'pending',
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_invoices_customer ON invoices(customer_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_period ON invoices(billing_period_start, billing_period_end);

-- 创建更新时间触发器
CREATE TRIGGER update_invoices_updated_at
    BEFORE UPDATE ON invoices
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE invoices IS '账单表';
COMMENT ON COLUMN invoices.status IS '状态: pending-待支付, paid-已支付, overdue-逾期, cancelled-已取消';
