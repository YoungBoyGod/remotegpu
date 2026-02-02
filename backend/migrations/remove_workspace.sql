-- 删除工作空间功能的数据库迁移脚本
-- 执行前请备份数据库！

-- 1. 移除 environments 表的 workspace_id 字段
ALTER TABLE environments DROP COLUMN IF EXISTS workspace_id;

-- 2. 移除 resource_quotas 表的 workspace_id 字段
-- 注意：需要先删除唯一索引
DROP INDEX IF EXISTS idx_user_workspace;
ALTER TABLE resource_quotas DROP COLUMN IF EXISTS workspace_id;
-- 重新创建用户级别的唯一索引
CREATE UNIQUE INDEX idx_user_quota ON resource_quotas(user_id) WHERE workspace_id IS NULL;

-- 3. 移除 datasets 表的 workspace_id 字段（如果存在）
ALTER TABLE datasets DROP COLUMN IF EXISTS workspace_id;

-- 4. 删除 workspace_members 表
DROP TABLE IF EXISTS workspace_members;

-- 5. 删除 workspaces 表
DROP TABLE IF EXISTS workspaces;

-- 迁移完成
