-- ============================================
-- 添加 SSH 连接信息字段
-- ============================================
-- 文件: 19_add_ssh_info.sql
-- 说明: 在 hosts 表中添加 SSH 连接相关字段
-- 执行顺序: 19
-- ============================================

-- 添加 SSH 连接主机地址（可与 ip_address/public_ip 不同）
ALTER TABLE hosts ADD COLUMN IF NOT EXISTS ssh_host VARCHAR(255);

-- 添加 SSH 用户名（部分环境可能已有该列，使用 IF NOT EXISTS）
ALTER TABLE hosts ADD COLUMN IF NOT EXISTS ssh_username VARCHAR(100) DEFAULT 'root';

-- 添加 SSH 密码（加密存储）
ALTER TABLE hosts ADD COLUMN IF NOT EXISTS ssh_password TEXT;

-- 添加 SSH 密钥（加密存储）
ALTER TABLE hosts ADD COLUMN IF NOT EXISTS ssh_key TEXT;

-- 添加注释
COMMENT ON COLUMN hosts.ssh_host IS 'SSH 连接主机地址，为空时使用 public_ip 或 ip_address';
COMMENT ON COLUMN hosts.ssh_username IS 'SSH 登录用户名，默认 root';
COMMENT ON COLUMN hosts.ssh_password IS 'SSH 登录密码（AES-256-GCM 加密存储）';
COMMENT ON COLUMN hosts.ssh_key IS 'SSH 私钥（AES-256-GCM 加密存储）';
