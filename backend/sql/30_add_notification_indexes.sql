-- ============================================
-- RemoteGPU notifications 表索引补全
-- ============================================
-- 文件: 30_add_notification_indexes.sql
-- 说明: 补全 notifications 表在 29_add_missing_indexes.sql 中遗漏的索引
-- ============================================

-- (customer_id, is_read, created_at DESC): ListByCustomerID() 分页查询覆盖索引
-- 覆盖 WHERE customer_id = ? [AND is_read = ?] ORDER BY created_at DESC 场景
CREATE INDEX IF NOT EXISTS idx_notifications_customer_read_created
    ON notifications(customer_id, is_read, created_at DESC);

-- type: 按通知类型筛选（system/environment/billing/alert）
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
