-- ============================================
-- 添加 Jupyter 和 VNC 连接信息字段
-- ============================================
-- 文件: 21_add_jupyter_vnc_info.sql
-- 说明: 在 hosts 表中添加 Jupyter 和 VNC 相关字段
-- 执行顺序: 21
-- ============================================

ALTER TABLE hosts ADD COLUMN IF NOT EXISTS jupyter_url VARCHAR(255);
ALTER TABLE hosts ADD COLUMN IF NOT EXISTS jupyter_token VARCHAR(255);
ALTER TABLE hosts ADD COLUMN IF NOT EXISTS vnc_url VARCHAR(255);
ALTER TABLE hosts ADD COLUMN IF NOT EXISTS vnc_password VARCHAR(255);

COMMENT ON COLUMN hosts.jupyter_url IS 'Jupyter Notebook 访问地址';
COMMENT ON COLUMN hosts.jupyter_token IS 'Jupyter Notebook 认证 Token';
COMMENT ON COLUMN hosts.vnc_url IS 'VNC 远程桌面访问地址';
COMMENT ON COLUMN hosts.vnc_password IS 'VNC 连接密码';
