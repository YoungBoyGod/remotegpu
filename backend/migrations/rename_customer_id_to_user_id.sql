-- 将 customer_id 重命名为 user_id 的数据库迁移脚本
-- 执行前请备份数据库！
-- 目的：统一使用 user_id 作为用户标识字段

-- ============================================
-- 第一部分：核心业务表
-- ============================================

-- 1. 修改 resource_quotas 表
-- 删除旧索引
DROP INDEX IF EXISTS idx_resource_quotas_customer;

-- 重命名字段
ALTER TABLE resource_quotas RENAME COLUMN customer_id TO user_id;

-- 创建新索引
CREATE INDEX idx_resource_quotas_user ON resource_quotas(user_id);

-- 2. 修改 environments 表
-- 删除旧索引
DROP INDEX IF EXISTS idx_environments_customer;

-- 重命名字段
ALTER TABLE environments RENAME COLUMN customer_id TO user_id;

-- 创建新索引
CREATE INDEX idx_environments_user ON environments(user_id);

-- ============================================
-- 第二部分：数据和权限表
-- ============================================

-- 3. 修改 datasets 表
DROP INDEX IF EXISTS idx_datasets_customer;
ALTER TABLE datasets RENAME COLUMN customer_id TO user_id;
CREATE INDEX idx_datasets_user ON datasets(user_id);

-- 4. 修改 images 表
DROP INDEX IF EXISTS idx_images_customer;
ALTER TABLE images RENAME COLUMN customer_id TO user_id;
CREATE INDEX idx_images_user ON images(user_id);

-- 5. 修改 api_keys 表
DROP INDEX IF EXISTS idx_api_keys_customer;
ALTER TABLE api_keys RENAME COLUMN customer_id TO user_id;
CREATE INDEX idx_api_keys_user ON api_keys(user_id);

-- ============================================
-- 第三部分：告警和计费表
-- ============================================

-- 6. 修改 alert_rules 表
ALTER TABLE alert_rules RENAME COLUMN customer_id TO user_id;

-- 7. 修改 alert_history 表
ALTER TABLE alert_history RENAME COLUMN customer_id TO user_id;

-- 8. 修改 webhooks 表
DROP INDEX IF EXISTS idx_webhooks_customer;
ALTER TABLE webhooks RENAME COLUMN customer_id TO user_id;
CREATE INDEX idx_webhooks_user ON webhooks(user_id);

-- 9. 修改 billing_accounts 表
DROP INDEX IF EXISTS idx_billing_accounts_customer;
ALTER TABLE billing_accounts RENAME COLUMN customer_id TO user_id;
CREATE INDEX idx_billing_accounts_user ON billing_accounts(user_id);

-- 10. 修改 invoices 表
DROP INDEX IF EXISTS idx_invoices_customer;
ALTER TABLE invoices RENAME COLUMN customer_id TO user_id;
CREATE INDEX idx_invoices_user ON invoices(user_id);

-- ============================================
-- 第四部分：训练、通知和问题管理表
-- ============================================

-- 11. 修改 training_jobs 表
DROP INDEX IF EXISTS idx_training_jobs_customer;
ALTER TABLE training_jobs RENAME COLUMN customer_id TO user_id;
CREATE INDEX idx_training_jobs_user ON training_jobs(user_id);

-- 12. 修改 inference_services 表
DROP INDEX IF EXISTS idx_inference_services_customer;
ALTER TABLE inference_services RENAME COLUMN customer_id TO user_id;
CREATE INDEX idx_inference_services_user ON inference_services(user_id);

-- 13. 修改 notifications 表
DROP INDEX IF EXISTS idx_notifications_customer;
ALTER TABLE notifications RENAME COLUMN customer_id TO user_id;
CREATE INDEX idx_notifications_user ON notifications(user_id);

-- 14. 修改 operation_logs 表
ALTER TABLE operation_logs RENAME COLUMN customer_id TO user_id;

-- 15. 修改 issues 表
DROP INDEX IF EXISTS idx_issues_customer;
ALTER TABLE issues RENAME COLUMN customer_id TO user_id;
CREATE INDEX idx_issues_user ON issues(user_id);

-- 16. 修改 requirements 表
DROP INDEX IF EXISTS idx_requirements_customer;
ALTER TABLE requirements RENAME COLUMN customer_id TO user_id;
CREATE INDEX idx_requirements_user ON requirements(user_id);

-- 17. 修改 feature_requests 表
DROP INDEX IF EXISTS idx_feature_requests_customer;
ALTER TABLE feature_requests RENAME COLUMN customer_id TO user_id;
CREATE INDEX idx_feature_requests_user ON feature_requests(user_id);

-- 18. 修改 user_favorites 表
DROP INDEX IF EXISTS idx_user_favorites_customer;
ALTER TABLE user_favorites RENAME COLUMN customer_id TO user_id;
CREATE INDEX idx_user_favorites_user ON user_favorites(user_id);

-- ============================================
-- 迁移完成
-- ============================================
