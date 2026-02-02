-- 扩展 port_mappings 表,添加协议和描述字段

-- 添加协议字段
ALTER TABLE port_mappings ADD COLUMN IF NOT EXISTS protocol VARCHAR(10) DEFAULT 'tcp';
COMMENT ON COLUMN port_mappings.protocol IS '协议类型: tcp/udp';

-- 添加描述字段
ALTER TABLE port_mappings ADD COLUMN IF NOT EXISTS description VARCHAR(256);
COMMENT ON COLUMN port_mappings.description IS '端口映射描述';

-- 更新 service_type 注释,添加更多服务类型
COMMENT ON COLUMN port_mappings.service_type IS '服务类型: ssh/rdp/jupyter/vnc/tensorboard/vscode/novnc/custom';
