-- ============================================
-- RemoteGPU 数据库初始化脚本
-- ============================================
-- 文件: 01_init_database.sql
-- 说明: 创建数据库和启用必要的扩展
-- 执行顺序: 1
-- ============================================

-- 创建数据库（如果不存在）
-- 注意：此命令需要在 postgres 数据库中执行
-- CREATE DATABASE remotegpu WITH ENCODING 'UTF8' LC_COLLATE='en_US.UTF-8' LC_CTYPE='en_US.UTF-8';

-- 连接到 remotegpu 数据库后执行以下命令

-- 启用 UUID 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 启用 pgcrypto 扩展（用于密码加密）
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 设置时区为 UTC
SET timezone = 'UTC';

-- 创建更新时间戳的触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION update_updated_at_column() IS '自动更新 updated_at 字段的触发器函数';
