-- ============================================
-- RemoteGPU 缺失索引补全
-- ============================================
-- 文件: 29_add_missing_indexes.sql
-- 说明: 根据 DAO 查询模式审查，补全缺失的单字段索引和复合索引
-- 优先级: 高（涉及核心业务表的高频查询）
-- ============================================

-- ============================================
-- 1. hosts 表：缺失的单字段索引
-- ============================================

-- ip_address: FindByIPAddress() 唯一性校验
CREATE INDEX IF NOT EXISTS idx_hosts_ip_address ON hosts(ip_address);

-- hostname: FindByHostname() 唯一性校验
CREATE INDEX IF NOT EXISTS idx_hosts_hostname ON hosts(hostname);

-- region: List() 管理端列表筛选
CREATE INDEX IF NOT EXISTS idx_hosts_region ON hosts(region);

-- needs_collect: ListNeedCollect() 采集任务查询
CREATE INDEX IF NOT EXISTS idx_hosts_needs_collect ON hosts(needs_collect) WHERE needs_collect = true;

-- last_heartbeat: HeartbeatMonitor 心跳超时检测排序
CREATE INDEX IF NOT EXISTS idx_hosts_last_heartbeat ON hosts(last_heartbeat);

-- ============================================
-- 2. tasks 表：完全缺失索引（高优先级）
-- ============================================

-- customer_id: ListByCustomerID() 客户查看自己的任务
CREATE INDEX IF NOT EXISTS idx_tasks_customer_id ON tasks(customer_id);

-- status: ListAll() 管理端任务列表筛选
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);

-- machine_id + status: ClaimTasks() Agent 认领任务（原子操作）
CREATE INDEX IF NOT EXISTS idx_tasks_machine_status ON tasks(machine_id, status);

-- ============================================
-- 3. allocations 表：缺失的复合索引（高优先级）
-- ============================================

-- (customer_id, status): FindAllActiveByCustomerID()
CREATE INDEX IF NOT EXISTS idx_allocations_customer_status ON allocations(customer_id, status);

-- (host_id, status): FindActiveByHostID()
CREATE INDEX IF NOT EXISTS idx_allocations_host_status ON allocations(host_id, status);

-- (host_id, customer_id, status): FindActiveByHostAndCustomer()
CREATE INDEX IF NOT EXISTS idx_allocations_host_customer_status ON allocations(host_id, customer_id, status);

-- ============================================
-- 4. dataset_mounts 表：缺失的复合索引（高优先级）
-- ============================================

-- (dataset_id, status): ListByDatasetID()
CREATE INDEX IF NOT EXISTS idx_dataset_mounts_dataset_status ON dataset_mounts(dataset_id, status);

-- (host_id, status): ListByHostID()、CountByHostID()
CREATE INDEX IF NOT EXISTS idx_dataset_mounts_host_status ON dataset_mounts(host_id, status);

-- (dataset_id, host_id, status): FindActiveMount()
CREATE INDEX IF NOT EXISTS idx_dataset_mounts_dataset_host_status ON dataset_mounts(dataset_id, host_id, status);

-- ============================================
-- 5. notifications 表：缺失的复合索引（中优先级）
-- ============================================

-- (customer_id, is_read): ListByCustomerID()、CountUnread()、MarkAllRead()
CREATE INDEX IF NOT EXISTS idx_notifications_customer_read ON notifications(customer_id, is_read);

-- ============================================
-- 6. customers 表：缺失的 role 索引（中优先级）
-- ============================================

-- role: 按角色筛选客户列表（GORM 实体定义了 index 标签但 SQL 迁移未创建）
-- 注意：customers 表使用 user_type 字段而非 role，此处为 user_type 补充覆盖
-- user_type 索引已存在于 03_users_and_permissions.sql，跳过

-- ============================================
-- 7. documents 表：外键 ON DELETE 策略修复
-- ============================================

-- 修复 uploaded_by 外键缺少 ON DELETE 策略
-- 先删除旧约束再重建（如果存在）
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'documents_uploaded_by_fkey'
        AND table_name = 'documents'
    ) THEN
        ALTER TABLE documents DROP CONSTRAINT documents_uploaded_by_fkey;
        ALTER TABLE documents ADD CONSTRAINT documents_uploaded_by_fkey
            FOREIGN KEY (uploaded_by) REFERENCES customers(id) ON DELETE SET NULL;
    END IF;
END $$;
