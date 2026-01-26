-- ============================================
-- RemoteGPU 数据和镜像管理表
-- ============================================
-- 文件: 06_data_and_images.sql
-- 说明: 创建数据集、模型、镜像相关表
-- 执行顺序: 6
-- ============================================

-- 数据集表
CREATE TABLE IF NOT EXISTS datasets (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    customer_id BIGINT NOT NULL,
    workspace_id BIGINT,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    storage_path VARCHAR(512) NOT NULL,
    storage_type VARCHAR(32) DEFAULT 'minio',
    total_size BIGINT DEFAULT 0,
    file_count INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'uploading',
    visibility VARCHAR(20) DEFAULT 'private',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_datasets_customer ON datasets(customer_id);
CREATE INDEX idx_datasets_workspace ON datasets(workspace_id);
CREATE INDEX idx_datasets_status ON datasets(status);
CREATE INDEX idx_datasets_visibility ON datasets(visibility);

-- 创建更新时间触发器
CREATE TRIGGER update_datasets_updated_at
    BEFORE UPDATE ON datasets
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE datasets IS '数据集表';
COMMENT ON COLUMN datasets.storage_type IS '存储类型: minio, s3, nfs';
COMMENT ON COLUMN datasets.status IS '状态: uploading-上传中, ready-就绪, error-错误';
COMMENT ON COLUMN datasets.visibility IS '可见性: private-私有, workspace-工作空间, public-公开';

-- 数据集版本表
CREATE TABLE IF NOT EXISTS dataset_versions (
    id BIGSERIAL PRIMARY KEY,
    dataset_id BIGINT NOT NULL,
    version VARCHAR(64) NOT NULL,
    description TEXT,
    storage_path VARCHAR(512) NOT NULL,
    size BIGINT DEFAULT 0,
    file_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(dataset_id, version)
);

-- 创建索引
CREATE INDEX idx_dataset_versions_dataset ON dataset_versions(dataset_id);

-- 添加注释
COMMENT ON TABLE dataset_versions IS '数据集版本表';
COMMENT ON COLUMN dataset_versions.version IS '版本号: v1.0, v1.1等';

-- 模型表
CREATE TABLE IF NOT EXISTS models (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    customer_id BIGINT NOT NULL,
    workspace_id BIGINT,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    framework VARCHAR(64),
    storage_path VARCHAR(512) NOT NULL,
    total_size BIGINT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'uploading',
    visibility VARCHAR(20) DEFAULT 'private',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_models_customer ON models(customer_id);
CREATE INDEX idx_models_workspace ON models(workspace_id);
CREATE INDEX idx_models_framework ON models(framework);

-- 创建更新时间触发器
CREATE TRIGGER update_models_updated_at
    BEFORE UPDATE ON models
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE models IS '模型表';
COMMENT ON COLUMN models.framework IS '框架: pytorch, tensorflow, onnx等';
COMMENT ON COLUMN models.status IS '状态: uploading-上传中, ready-就绪, error-错误';

-- 模型版本表
CREATE TABLE IF NOT EXISTS model_versions (
    id BIGSERIAL PRIMARY KEY,
    model_id BIGINT NOT NULL,
    version VARCHAR(64) NOT NULL,
    description TEXT,
    storage_path VARCHAR(512) NOT NULL,
    size BIGINT DEFAULT 0,
    metrics JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(model_id, version)
);

-- 创建索引
CREATE INDEX idx_model_versions_model ON model_versions(model_id);

-- 添加注释
COMMENT ON TABLE model_versions IS '模型版本表';
COMMENT ON COLUMN model_versions.metrics IS '模型指标: accuracy, loss等';

-- 镜像表
CREATE TABLE IF NOT EXISTS images (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) UNIQUE NOT NULL,
    display_name VARCHAR(256),
    description TEXT,
    category VARCHAR(64),
    framework VARCHAR(64),
    framework_version VARCHAR(64),
    cuda_version VARCHAR(32),
    python_version VARCHAR(32),
    is_official BOOLEAN DEFAULT false,
    size BIGINT,
    registry_url VARCHAR(512),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_images_category ON images(category);
CREATE INDEX idx_images_framework ON images(framework);
CREATE INDEX idx_images_is_official ON images(is_official);
CREATE INDEX idx_images_status ON images(status);

-- 创建更新时间触发器
CREATE TRIGGER update_images_updated_at
    BEFORE UPDATE ON images
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE images IS '镜像表';
COMMENT ON COLUMN images.category IS '分类: base, pytorch, tensorflow, custom等';
COMMENT ON COLUMN images.status IS '状态: active-活跃, deprecated-已弃用';
