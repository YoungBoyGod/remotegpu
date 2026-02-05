-- ============================================
-- RemoteGPU task process id
-- ============================================
-- File: 15_task_process_id.sql
-- Description: add process_id column for agent stop integration
-- Execution order: 15
-- ============================================

ALTER TABLE IF EXISTS tasks
    ADD COLUMN IF NOT EXISTS process_id INT DEFAULT 0;

COMMENT ON COLUMN tasks.process_id IS 'Agent process id for stop integration';
