-- ============================================
-- RemoteGPU 训练和推理表
-- ============================================
-- 文件: 09_training_and_inference.sql
-- 说明: 创建训练任务、推理服务相关表
-- 执行顺序: 9
-- ============================================

-- 训练任务表
CREATE TABLE IF NOT EXISTS training_jobs (
    id VARCHAR(64) PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    workspace_id BIGINT,
    env_id VARCHAR(64),
    name VARCHAR(256) NOT NULL,
    description TEXT,
    framework VARCHAR(64),
    script_path VARCHAR(512),
    dataset_id BIGINT,
    model_id BIGINT,
    status VARCHAR(20) DEFAULT 'pending',
    priority INT DEFAULT 0,
    cpu INT,
    memory BIGINT,
    gpu INT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_training_jobs_customer ON training_jobs(customer_id);
CREATE INDEX idx_training_jobs_workspace ON training_jobs(workspace_id);
CREATE INDEX idx_training_jobs_status ON training_jobs(status);
CREATE INDEX idx_training_jobs_created_at ON training_jobs(created_at DESC);

-- 创建更新时间触发器
CREATE TRIGGER update_training_jobs_updated_at
    BEFORE UPDATE ON training_jobs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE training_jobs IS '训练任务表';
COMMENT ON COLUMN training_jobs.status IS '状态: pending-待运行, running-运行中, completed-已完成, failed-失败, cancelled-已取消';
COMMENT ON COLUMN training_jobs.priority IS '优先级: 数字越大优先级越高';

-- 推理服务表
CREATE TABLE IF NOT EXISTS inference_services (
    id VARCHAR(64) PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    workspace_id BIGINT,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    model_id BIGINT,
    model_version VARCHAR(64),
    framework VARCHAR(64),
    status VARCHAR(20) DEFAULT 'creating',
    endpoint_url VARCHAR(512),
    cpu INT,
    memory BIGINT,
    gpu INT,
    replicas INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_inference_services_customer ON inference_services(customer_id);
CREATE INDEX idx_inference_services_workspace ON inference_services(workspace_id);
CREATE INDEX idx_inference_services_status ON inference_services(status);

-- 创建更新时间触发器
CREATE TRIGGER update_inference_services_updated_at
    BEFORE UPDATE ON inference_services
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE inference_services IS '推理服务表';
COMMENT ON COLUMN inference_services.status IS '状态: creating-创建中, running-运行中, stopped-已停止, error-错误';
COMMENT ON COLUMN inference_services.replicas IS '副本数量';
