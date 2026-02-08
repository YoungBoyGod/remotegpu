-- ============================================
-- 为 environments 表添加 workspace_id 字段
-- ============================================
-- 文件: 31_add_environment_workspace_id.sql
-- 说明: 环境关联工作空间，支持按工作空间筛选环境
-- 执行顺序: 31
-- ============================================

ALTER TABLE environments ADD COLUMN IF NOT EXISTS workspace_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_environments_workspace ON environments(workspace_id);

COMMENT ON COLUMN environments.workspace_id IS '所属工作空间ID（可选）';
