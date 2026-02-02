-- 扩展 port_mappings 表,添加公网访问和域名字段

-- 添加公网端口字段
ALTER TABLE port_mappings ADD COLUMN IF NOT EXISTS public_port INT;
COMMENT ON COLUMN port_mappings.public_port IS '公网端口(防火墙映射后的端口)';

-- 添加内网访问地址
ALTER TABLE port_mappings ADD COLUMN IF NOT EXISTS internal_access_url VARCHAR(512);
COMMENT ON COLUMN port_mappings.internal_access_url IS '内网访问地址 (如 192.168.1.100:22001)';

-- 添加公网访问地址
ALTER TABLE port_mappings ADD COLUMN IF NOT EXISTS public_access_url VARCHAR(512);
COMMENT ON COLUMN port_mappings.public_access_url IS '公网访问地址 (如 ssh-env123.example.com:22001)';

-- 添加公网域名
ALTER TABLE port_mappings ADD COLUMN IF NOT EXISTS public_domain VARCHAR(256);
COMMENT ON COLUMN port_mappings.public_domain IS '公网域名 (如 ssh-env123.example.com)';

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_port_mappings_public_domain ON port_mappings(public_domain);
