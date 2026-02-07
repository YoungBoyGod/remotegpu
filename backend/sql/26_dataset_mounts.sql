-- ============================================
-- RemoteGPU 数据集挂载表
-- ============================================
-- 文件: 26_dataset_mounts.sql
-- 说明: 创建数据集挂载管理表，记录数据集到机器的挂载关系
-- ============================================

CREATE TABLE IF NOT EXISTS dataset_mounts (
    id BIGSERIAL PRIMARY KEY,
    dataset_id BIGINT NOT NULL,
    host_id VARCHAR(64) NOT NULL,
    mount_path VARCHAR(256) NOT NULL,
    read_only BOOLEAN DEFAULT true,
    status VARCHAR(20) DEFAULT 'mounting',
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_dataset_mounts_dataset ON dataset_mounts(dataset_id);
CREATE INDEX IF NOT EXISTS idx_dataset_mounts_host ON dataset_mounts(host_id);
CREATE INDEX IF NOT EXISTS idx_dataset_mounts_status ON dataset_mounts(status);

-- 创建更新时间触发器
CREATE TRIGGER update_dataset_mounts_updated_at
    BEFORE UPDATE ON dataset_mounts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE dataset_mounts IS '数据集挂载表';
COMMENT ON COLUMN dataset_mounts.dataset_id IS '关联的数据集ID';
COMMENT ON COLUMN dataset_mounts.host_id IS '挂载目标机器ID';
COMMENT ON COLUMN dataset_mounts.mount_path IS '挂载路径（绝对路径）';
COMMENT ON COLUMN dataset_mounts.read_only IS '是否只读挂载';
COMMENT ON COLUMN dataset_mounts.status IS '状态: mounting-挂载中, mounted-已挂载, error-错误, unmounted-已卸载';
COMMENT ON COLUMN dataset_mounts.error_message IS '错误信息';
