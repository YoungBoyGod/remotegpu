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

    -- 镜像和框架
    image VARCHAR(256) NOT NULL,
    framework VARCHAR(64),

    -- 挂载配置 (JSONB格式: [1, 2, 3] 最多3个)
    mounted_datasets JSONB,
    mounted_models JSONB,

    -- 执行配置
    command TEXT NOT NULL,
    env_vars JSONB,
    output_path VARCHAR(512) DEFAULT '/gemini/output/',

    -- 分布式训练配置
    node_count INT DEFAULT 1,
    distributed_config JSONB,

    -- 任务状态
    status VARCHAR(20) DEFAULT 'pending',
    priority INT DEFAULT 0,

    -- 资源配置
    cpu INT,
    memory BIGINT,
    gpu INT,
    gpu_memory BIGINT,

    -- 时间戳
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_training_jobs_customer ON training_jobs(customer_id);
CREATE INDEX idx_training_jobs_workspace ON training_jobs(workspace_id);
CREATE INDEX idx_training_jobs_status ON training_jobs(status);
CREATE INDEX idx_training_jobs_created_at ON training_jobs(created_at DESC);
CREATE INDEX idx_training_jobs_deleted_at ON training_jobs(deleted_at);

-- 为JSONB字段创建GIN索引
CREATE INDEX idx_training_jobs_mounted_datasets ON training_jobs USING GIN (mounted_datasets);
CREATE INDEX idx_training_jobs_mounted_models ON training_jobs USING GIN (mounted_models);
CREATE INDEX idx_training_jobs_env_vars ON training_jobs USING GIN (env_vars);
CREATE INDEX idx_training_jobs_distributed_config ON training_jobs USING GIN (distributed_config);

-- 创建更新时间触发器
CREATE TRIGGER update_training_jobs_updated_at
    BEFORE UPDATE ON training_jobs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE training_jobs IS '训练任务表';
COMMENT ON COLUMN training_jobs.status IS '状态: pending-待运行, queued-排队中, running-运行中, completed-已完成, failed-失败, cancelled-已取消';
COMMENT ON COLUMN training_jobs.priority IS '优先级: 数字越大优先级越高';
COMMENT ON COLUMN training_jobs.mounted_datasets IS '挂载的数据集ID列表(JSONB): [1, 2, 3] 最多3个';
COMMENT ON COLUMN training_jobs.mounted_models IS '挂载的模型ID列表(JSONB): [1, 2] 最多3个';
COMMENT ON COLUMN training_jobs.command IS '启动命令';
COMMENT ON COLUMN training_jobs.env_vars IS '环境变量(JSONB): {"KEY": "value"}';
COMMENT ON COLUMN training_jobs.node_count IS '分布式训练节点数量';
COMMENT ON COLUMN training_jobs.distributed_config IS '分布式训练配置(JSONB): {"framework": "pytorch", "backend": "nccl"}';
COMMENT ON COLUMN training_jobs.gpu_memory IS 'GPU显存(字节)';
COMMENT ON COLUMN training_jobs.deleted_at IS '软删除时间';

-- 推理服务表
CREATE TABLE IF NOT EXISTS inference_services (
    id VARCHAR(64) PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    workspace_id BIGINT,
    name VARCHAR(256) NOT NULL,
    description TEXT,

    -- 模型和镜像
    model_id BIGINT,
    model_version VARCHAR(64),
    image VARCHAR(256) NOT NULL,
    framework VARCHAR(64),

    -- 服务状态
    status VARCHAR(20) DEFAULT 'creating',
    endpoint_url VARCHAR(512),

    -- 副本配置
    replicas INT DEFAULT 1,
    min_replicas INT DEFAULT 1,
    max_replicas INT DEFAULT 10,

    -- 自动扩缩容
    autoscaling_enabled BOOLEAN DEFAULT false,
    autoscaling_config JSONB,

    -- 版本管理（支持回滚）
    version VARCHAR(64),
    previous_version VARCHAR(64),

    -- 健康检查
    health_check_path VARCHAR(256),
    health_check_interval INT DEFAULT 30,

    -- 环境配置
    env_vars JSONB,

    -- 资源配置
    cpu INT,
    memory BIGINT,
    gpu INT,
    gpu_memory BIGINT,

    -- 时间戳
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_inference_services_customer ON inference_services(customer_id);
CREATE INDEX idx_inference_services_workspace ON inference_services(workspace_id);
CREATE INDEX idx_inference_services_status ON inference_services(status);
CREATE INDEX idx_inference_services_deleted_at ON inference_services(deleted_at);

-- 为JSONB字段创建GIN索引
CREATE INDEX idx_inference_services_autoscaling_config ON inference_services USING GIN (autoscaling_config);
CREATE INDEX idx_inference_services_env_vars ON inference_services USING GIN (env_vars);

-- 创建更新时间触发器
CREATE TRIGGER update_inference_services_updated_at
    BEFORE UPDATE ON inference_services
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE inference_services IS '推理服务表';
COMMENT ON COLUMN inference_services.status IS '状态: creating-创建中, running-运行中, stopped-已停止, error-错误';
COMMENT ON COLUMN inference_services.replicas IS '当前副本数量';
COMMENT ON COLUMN inference_services.min_replicas IS '最小副本数量';
COMMENT ON COLUMN inference_services.max_replicas IS '最大副本数量';
COMMENT ON COLUMN inference_services.autoscaling_enabled IS '是否启用自动扩缩容';
COMMENT ON COLUMN inference_services.autoscaling_config IS '扩缩容配置(JSONB): {"target_cpu": 80, "target_qps": 1000}';
COMMENT ON COLUMN inference_services.version IS '当前版本';
COMMENT ON COLUMN inference_services.previous_version IS '上一个版本（用于回滚）';
COMMENT ON COLUMN inference_services.health_check_path IS '健康检查路径';
COMMENT ON COLUMN inference_services.health_check_interval IS '健康检查间隔(秒)';
COMMENT ON COLUMN inference_services.env_vars IS '环境变量(JSONB): {"KEY": "value"}';
COMMENT ON COLUMN inference_services.gpu_memory IS 'GPU显存(字节)';
COMMENT ON COLUMN inference_services.deleted_at IS '软删除时间';
