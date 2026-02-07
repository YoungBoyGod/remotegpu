-- 客户配额字段
ALTER TABLE customers ADD COLUMN IF NOT EXISTS quota_gpu INTEGER DEFAULT 0;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS quota_storage BIGINT DEFAULT 0;

COMMENT ON COLUMN customers.quota_gpu IS 'GPU 数量配额，0 表示不限制';
COMMENT ON COLUMN customers.quota_storage IS '存储容量配额（MB），0 表示不限制';
