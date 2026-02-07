-- 为 audit_logs 表添加 detail 字段，存储请求参数摘要
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS detail JSONB;

COMMENT ON COLUMN audit_logs.detail IS '请求参数摘要（JSON 格式，敏感字段已脱敏）';
