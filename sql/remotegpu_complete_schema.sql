-- RemoteGPU Database Schema (Consolidated & Aligned)
-- Version: 2.1 (Aligned with V2.0 Requirements & Existing Schema)
-- Date: 2026-02-03
-- Database: PostgreSQL 14+

-- 1. Enable Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ==========================================
-- 2. System Configurations
-- ==========================================
CREATE TABLE IF NOT EXISTS system_configs (
    id BIGSERIAL PRIMARY KEY,
    config_key VARCHAR(128) NOT NULL UNIQUE,
    config_value TEXT NOT NULL,
    config_type VARCHAR(32) DEFAULT 'string',
    description TEXT,
    is_public BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================
-- 3. Users & Tenants
-- ==========================================
CREATE TABLE IF NOT EXISTS customers (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    username VARCHAR(64) NOT NULL UNIQUE,
    email VARCHAR(128) NOT NULL UNIQUE,
    password_hash VARCHAR(256) NOT NULL,
    
    -- Profile
    display_name VARCHAR(128),
    full_name VARCHAR(256),
    company VARCHAR(256),
    phone VARCHAR(32),
    avatar_url VARCHAR(512),
    
    -- Role & Type (Aligned with existing schema + V2 requirements)
    role VARCHAR(32) DEFAULT 'customer_owner', -- admin, customer_owner, customer_member
    user_type VARCHAR(32) DEFAULT 'external', -- admin (internal), external (customer)
    account_type VARCHAR(32) DEFAULT 'individual', -- individual, enterprise
    
    -- Status
    status VARCHAR(32) DEFAULT 'active', -- active, suspended, deleted
    email_verified BOOLEAN DEFAULT FALSE,
    phone_verified BOOLEAN DEFAULT FALSE,
    
    -- Billing
    balance DECIMAL(10,4) DEFAULT 0.00,
    currency VARCHAR(10) DEFAULT 'CNY',
    
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_customers_username ON customers(username);
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_role ON customers(role);

-- SSH Keys
CREATE TABLE IF NOT EXISTS ssh_keys (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customers(id),
    name VARCHAR(64) NOT NULL,
    fingerprint VARCHAR(128),
    public_key TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Workspaces (Optional: for team isolation within a customer account)
CREATE TABLE IF NOT EXISTS workspaces (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    owner_id BIGINT NOT NULL REFERENCES customers(id),
    name VARCHAR(128) NOT NULL,
    description TEXT,
    type VARCHAR(32) DEFAULT 'personal', -- personal, team
    status VARCHAR(32) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================
-- 4. Resources: Machines & GPUs
-- ==========================================
CREATE TABLE IF NOT EXISTS hosts (
    id VARCHAR(64) PRIMARY KEY, -- Machine ID
    name VARCHAR(128) NOT NULL,
    hostname VARCHAR(256),
    region VARCHAR(64) DEFAULT 'default',
    
    -- Network
    ip_address VARCHAR(64) NOT NULL, -- Internal IP
    public_ip VARCHAR(64),           -- External IP
    ssh_port INT DEFAULT 22,
    agent_port INT DEFAULT 8080,
    
    -- Specs
    os_type VARCHAR(20) DEFAULT 'linux',
    os_version VARCHAR(64),
    cpu_info VARCHAR(256),
    total_cpu INT NOT NULL,
    total_memory_gb BIGINT NOT NULL,
    total_disk_gb BIGINT,
    
    -- Status
    status VARCHAR(20) DEFAULT 'offline', -- offline, idle, allocated, maintenance
    health_status VARCHAR(20) DEFAULT 'unknown', -- healthy, degraded, error
    deployment_mode VARCHAR(20) DEFAULT 'traditional', -- traditional, kubernetes
    
    last_heartbeat TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gpus (
    id BIGSERIAL PRIMARY KEY,
    host_id VARCHAR(64) NOT NULL REFERENCES hosts(id),
    index INT NOT NULL,
    uuid VARCHAR(128) UNIQUE,
    
    name VARCHAR(128) NOT NULL,
    memory_total_mb INT NOT NULL,
    brand VARCHAR(64),
    
    -- Status
    status VARCHAR(20) DEFAULT 'available', -- available, allocated, error
    health_status VARCHAR(20) DEFAULT 'healthy',
    allocated_to VARCHAR(64), -- Can reference allocation_id or container_id
    
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(host_id, index)
);

-- ==========================================
-- 5. Allocations (Core Business)
-- ==========================================
CREATE TABLE IF NOT EXISTS allocations (
    id VARCHAR(64) PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customers(id),
    host_id VARCHAR(64) NOT NULL REFERENCES hosts(id),
    workspace_id BIGINT REFERENCES workspaces(id),
    
    -- Time
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    actual_end_time TIMESTAMP WITH TIME ZONE,
    
    -- Status
    status VARCHAR(32) DEFAULT 'active', -- active, expired, reclaimed, pending
    
    -- Metadata
    remark TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_allocations_customer ON allocations(customer_id);
CREATE INDEX idx_allocations_status ON allocations(status);

-- ==========================================
-- 6. Images & Datasets
-- ==========================================
CREATE TABLE IF NOT EXISTS images (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL UNIQUE,
    display_name VARCHAR(256),
    description TEXT,
    
    category VARCHAR(64),
    framework VARCHAR(64),
    cuda_version VARCHAR(32),
    
    registry_url VARCHAR(512),
    is_official BOOLEAN DEFAULT false,
    customer_id BIGINT REFERENCES customers(id), -- If private image
    
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS datasets (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    customer_id BIGINT NOT NULL REFERENCES customers(id),
    workspace_id BIGINT REFERENCES workspaces(id),
    
    name VARCHAR(256) NOT NULL,
    description TEXT,
    
    storage_path VARCHAR(512) NOT NULL,
    storage_type VARCHAR(32) DEFAULT 'minio',
    total_size BIGINT DEFAULT 0,
    file_count INT DEFAULT 0,
    
    status VARCHAR(20) DEFAULT 'uploading', -- uploading, ready, error
    visibility VARCHAR(20) DEFAULT 'private',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS dataset_mounts (
    id BIGSERIAL PRIMARY KEY,
    dataset_id BIGINT NOT NULL REFERENCES datasets(id),
    host_id VARCHAR(64) NOT NULL REFERENCES hosts(id),
    mount_path VARCHAR(256) NOT NULL,
    read_only BOOLEAN DEFAULT true,
    
    status VARCHAR(20) DEFAULT 'mounting',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================
-- 7. Tasks (Training/Inference)
-- ==========================================
CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(64) PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customers(id),
    host_id VARCHAR(64) REFERENCES hosts(id),
    
    name VARCHAR(256) NOT NULL,
    type VARCHAR(32) NOT NULL, -- training, inference
    
    image_id BIGINT REFERENCES images(id),
    command TEXT NOT NULL,
    env_vars JSONB,
    
    status VARCHAR(20) DEFAULT 'queued',
    exit_code INT,
    error_msg TEXT,
    
    started_at TIMESTAMP WITH TIME ZONE,
    finished_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================
-- 8. Ops: Monitoring & Audit
-- ==========================================
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT,
    username VARCHAR(128),
    ip_address VARCHAR(64),
    method VARCHAR(10),
    path VARCHAR(512),
    
    action VARCHAR(128) NOT NULL,
    resource_type VARCHAR(64),
    resource_id VARCHAR(128),
    detail JSONB,
    
    status_code INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS alert_rules (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    metric_type VARCHAR(64) NOT NULL,
    threshold FLOAT NOT NULL,
    condition VARCHAR(10) NOT NULL,
    severity VARCHAR(20) DEFAULT 'warning',
    
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS active_alerts (
    id BIGSERIAL PRIMARY KEY,
    rule_id BIGINT REFERENCES alert_rules(id),
    host_id VARCHAR(64) NOT NULL REFERENCES hosts(id),
    
    value FLOAT NOT NULL,
    message TEXT,
    
    triggered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    acknowledged BOOLEAN DEFAULT false
);

-- Default Admin (Password: admin123)
INSERT INTO customers (username, email, password_hash, role, display_name, user_type) 
VALUES ('admin', 'admin@remotegpu.com', '$2a$10$NotRealHashJustExamplePlaceholder', 'admin', 'System Administrator', 'admin')
ON CONFLICT (username) DO NOTHING;