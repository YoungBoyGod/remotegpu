-- ============================================
-- 拆分机器状态为设备状态和分配状态
-- ============================================
-- 文件: 20_split_machine_status.sql
-- 说明: 将 hosts 表的 status 字段拆分为 device_status 和 allocation_status
-- 执行顺序: 20
-- ============================================

-- 添加设备状态字段（表示机器是否在线）
ALTER TABLE hosts ADD COLUMN IF NOT EXISTS device_status VARCHAR(20) DEFAULT 'offline';

-- 添加分配状态字段（表示机器的使用状态）
ALTER TABLE hosts ADD COLUMN IF NOT EXISTS allocation_status VARCHAR(20) DEFAULT 'idle';

-- 从现有 status 字段迁移数据
UPDATE hosts SET
    device_status = CASE
        WHEN status = 'offline' THEN 'offline'
        ELSE 'online'
    END,
    allocation_status = CASE
        WHEN status = 'allocated' THEN 'allocated'
        WHEN status = 'maintenance' THEN 'maintenance'
        ELSE 'idle'
    END;

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_hosts_device_status ON hosts(device_status);
CREATE INDEX IF NOT EXISTS idx_hosts_allocation_status ON hosts(allocation_status);

-- 添加注释
COMMENT ON COLUMN hosts.device_status IS '设备状态: online-在线, offline-离线';
COMMENT ON COLUMN hosts.allocation_status IS '分配状态: idle-空闲, allocated-已分配, maintenance-维护中';
