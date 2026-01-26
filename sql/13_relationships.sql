-- ============================================
-- RemoteGPU 关联关系表
-- ============================================
-- 文件: 13_relationships.sql
-- 说明: 创建数据集使用、制品等关联关系表
-- 执行顺序: 13
-- ============================================

-- 数据集使用记录表（环境挂载数据集）
CREATE TABLE IF NOT EXISTS dataset_usage (
    id BIGSERIAL PRIMARY KEY,
    env_id VARCHAR(64) NOT NULL,
    dataset_id BIGINT NOT NULL,
    mount_path VARCHAR(512) NOT NULL,
    readonly BOOLEAN DEFAULT true,
    mounted_at TIMESTAMP DEFAULT NOW(),
    unmounted_at TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_dataset_usage_env ON dataset_usage(env_id);
CREATE INDEX idx_dataset_usage_dataset ON dataset_usage(dataset_id);

-- 添加注释
COMMENT ON TABLE dataset_usage IS '数据集使用记录表';
COMMENT ON COLUMN dataset_usage.mount_path IS '挂载路径';
COMMENT ON COLUMN dataset_usage.readonly IS '是否只读';

-- 制品表
CREATE TABLE IF NOT EXISTS artifacts (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    workspace_id BIGINT,
    name VARCHAR(256) NOT NULL,
    type VARCHAR(64) NOT NULL,
    storage_path VARCHAR(512) NOT NULL,
    size BIGINT DEFAULT 0,
    source_type VARCHAR(32),
    source_id VARCHAR(128),
    created_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_artifacts_customer ON artifacts(customer_id);
CREATE INDEX idx_artifacts_workspace ON artifacts(workspace_id);
CREATE INDEX idx_artifacts_type ON artifacts(type);
CREATE INDEX idx_artifacts_source ON artifacts(source_type, source_id);

-- 添加注释
COMMENT ON TABLE artifacts IS '制品表';
COMMENT ON COLUMN artifacts.type IS '类型: model-模型, checkpoint-检查点, log-日志, report-报告';
COMMENT ON COLUMN artifacts.source_type IS '来源类型: training_job, inference_service, environment';
