-- 添加环境扩展字段
-- 用于支持多种部署模式、存储类型、桌面环境等

-- 添加部署模式字段
ALTER TABLE environments ADD COLUMN IF NOT EXISTS deployment_mode VARCHAR(32) DEFAULT 'k8s_pod';
COMMENT ON COLUMN environments.deployment_mode IS '部署模式: k8s_pod, k8s_stateful, docker_local, vm, bare_metal';

-- 添加环境类型字段
ALTER TABLE environments ADD COLUMN IF NOT EXISTS environment_type VARCHAR(32) DEFAULT 'ide';
COMMENT ON COLUMN environments.environment_type IS '环境类型: ide, terminal, desktop, training, inference, data_process';

-- 添加 GPU 模式字段
ALTER TABLE environments ADD COLUMN IF NOT EXISTS gpu_mode VARCHAR(32) DEFAULT 'exclusive';
COMMENT ON COLUMN environments.gpu_mode IS 'GPU 分配模式: exclusive, vgpu, mig, shared, none';

-- 添加存储类型字段
ALTER TABLE environments ADD COLUMN IF NOT EXISTS storage_type VARCHAR(32) DEFAULT 'local';
COMMENT ON COLUMN environments.storage_type IS '存储类型: local, nfs, ceph, s3, pvc, juicefs';

-- 添加存储配置字段
ALTER TABLE environments ADD COLUMN IF NOT EXISTS storage_config JSONB;
COMMENT ON COLUMN environments.storage_config IS '存储配置(JSON): NFS/S3/JuiceFS 等配置信息';

-- 添加生命周期策略字段
ALTER TABLE environments ADD COLUMN IF NOT EXISTS lifecycle_policy VARCHAR(32) DEFAULT 'persistent';
COMMENT ON COLUMN environments.lifecycle_policy IS '生命周期策略: persistent, ephemeral, scheduled, on_demand';

-- 添加网络配置字段
ALTER TABLE environments ADD COLUMN IF NOT EXISTS network_config JSONB;
COMMENT ON COLUMN environments.network_config IS '网络配置(JSON): 端口映射、域名等配置';

-- 添加 Jumpserver 标志字段
ALTER TABLE environments ADD COLUMN IF NOT EXISTS use_jumpserver BOOLEAN DEFAULT false;
COMMENT ON COLUMN environments.use_jumpserver IS '是否使用 Jumpserver 堡垒机';

-- 添加 VNC 端口字段
ALTER TABLE environments ADD COLUMN IF NOT EXISTS vnc_port INTEGER;
COMMENT ON COLUMN environments.vnc_port IS 'VNC 端口';

-- 添加 VNC 密码字段
ALTER TABLE environments ADD COLUMN IF NOT EXISTS vnc_password VARCHAR(128);
COMMENT ON COLUMN environments.vnc_password IS 'VNC 密码';

-- 添加额外端口映射字段
ALTER TABLE environments ADD COLUMN IF NOT EXISTS additional_ports JSONB;
COMMENT ON COLUMN environments.additional_ports IS '额外端口映射(JSON): TensorBoard, VSCode, 自定义服务等';

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_environments_deployment_mode ON environments(deployment_mode);
CREATE INDEX IF NOT EXISTS idx_environments_environment_type ON environments(environment_type);
CREATE INDEX IF NOT EXISTS idx_environments_storage_type ON environments(storage_type);
CREATE INDEX IF NOT EXISTS idx_environments_lifecycle_policy ON environments(lifecycle_policy);
