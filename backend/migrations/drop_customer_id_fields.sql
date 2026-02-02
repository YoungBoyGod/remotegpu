-- 删除 customer_id 字段的数据库迁移脚本
-- 执行前请备份数据库！
-- 目的：删除旧的 customer_id 字段,统一使用 user_id

-- ============================================
-- 第一部分：核心业务表
-- ============================================

-- 1. 修改 resource_quotas 表
-- 删除外键约束
ALTER TABLE resource_quotas DROP CONSTRAINT IF EXISTS fk_resource_quotas_customer;

-- 删除唯一索引
DROP INDEX IF EXISTS idx_customer_workspace;

-- 删除普通索引
DROP INDEX IF EXISTS idx_resource_quotas_customer;

-- 删除字段
ALTER TABLE resource_quotas DROP COLUMN IF EXISTS customer_id;

-- 2. 修改 environments 表
-- 删除外键约束
ALTER TABLE environments DROP CONSTRAINT IF EXISTS fk_environments_customer;

-- 删除索引
DROP INDEX IF EXISTS idx_environments_customer;

-- 删除字段
ALTER TABLE environments DROP COLUMN IF EXISTS customer_id;
