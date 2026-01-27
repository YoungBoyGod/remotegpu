-- ============================================
-- RemoteGPU 用户与权限表
-- ============================================
-- 文件: 03_users_and_permissions.sql
-- 说明: 创建用户、工作空间、权限相关表
-- 执行顺序: 3
-- ============================================

-- 客户表（用户表）
CREATE TABLE IF NOT EXISTS customers (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    username VARCHAR(64) UNIQUE NOT NULL,
    email VARCHAR(128) UNIQUE NOT NULL,
    password_hash VARCHAR(256) NOT NULL,
    display_name VARCHAR(128),
    avatar_url VARCHAR(512),
    phone VARCHAR(32),
    full_name VARCHAR(256),
    company VARCHAR(256),
    user_type VARCHAR(32) DEFAULT 'external',
    account_type VARCHAR(32) DEFAULT 'individual',
    status VARCHAR(32) DEFAULT 'active',
    email_verified BOOLEAN DEFAULT FALSE,
    phone_verified BOOLEAN DEFAULT FALSE,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_customers_username ON customers(username);
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_status ON customers(status);
CREATE INDEX idx_customers_user_type ON customers(user_type);
CREATE INDEX idx_customers_account_type ON customers(account_type);

-- 创建更新时间触发器
CREATE TRIGGER update_customers_updated_at
    BEFORE UPDATE ON customers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE customers IS '客户表（用户表）';
COMMENT ON COLUMN customers.user_type IS '用户类型: admin-管理人员, internal-内部用户, external-外部用户';
COMMENT ON COLUMN customers.account_type IS '账户类型: individual-个人, enterprise-企业';
COMMENT ON COLUMN customers.status IS '状态: active-活跃, suspended-暂停, deleted-已删除';

-- 工作空间表
CREATE TABLE IF NOT EXISTS workspaces (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    owner_id BIGINT NOT NULL,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    type VARCHAR(32) DEFAULT 'personal',
    member_count INT DEFAULT 1,
    status VARCHAR(32) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_workspaces_owner ON workspaces(owner_id);
CREATE INDEX idx_workspaces_type ON workspaces(type);
CREATE INDEX idx_workspaces_status ON workspaces(status);

-- 创建更新时间触发器
CREATE TRIGGER update_workspaces_updated_at
    BEFORE UPDATE ON workspaces
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE workspaces IS '工作空间表';
COMMENT ON COLUMN workspaces.type IS '类型: personal-个人, team-团队, enterprise-企业';
COMMENT ON COLUMN workspaces.status IS '状态: active-活跃, archived-归档';

-- 工作空间成员表
CREATE TABLE IF NOT EXISTS workspace_members (
    id BIGSERIAL PRIMARY KEY,
    workspace_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    role VARCHAR(32) DEFAULT 'member',
    status VARCHAR(32) DEFAULT 'active',
    joined_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(workspace_id, customer_id)
);

-- 创建索引
CREATE INDEX idx_workspace_members_workspace ON workspace_members(workspace_id);
CREATE INDEX idx_workspace_members_customer ON workspace_members(customer_id);
CREATE INDEX idx_workspace_members_role ON workspace_members(role);

-- 添加注释
COMMENT ON TABLE workspace_members IS '工作空间成员表';
COMMENT ON COLUMN workspace_members.role IS '角色: owner-所有者, admin-管理员, member-成员, viewer-查看者';
COMMENT ON COLUMN workspace_members.status IS '状态: active-活跃, invited-已邀请, suspended-暂停';

-- 资源配额表
CREATE TABLE IF NOT EXISTS resource_quotas (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    workspace_id BIGINT,
    quota_level VARCHAR(32) DEFAULT 'free',
    cpu_quota INT DEFAULT 4,
    memory_quota INT DEFAULT 8192,
    gpu_quota INT DEFAULT 0,
    storage_quota BIGINT DEFAULT 10737418240,
    environment_quota INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_resource_quotas_customer ON resource_quotas(customer_id);
CREATE INDEX idx_resource_quotas_workspace ON resource_quotas(workspace_id);

-- 创建更新时间触发器
CREATE TRIGGER update_resource_quotas_updated_at
    BEFORE UPDATE ON resource_quotas
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE resource_quotas IS '资源配额表';
COMMENT ON COLUMN resource_quotas.quota_level IS '配额级别: free, basic, pro, enterprise';
COMMENT ON COLUMN resource_quotas.cpu_quota IS 'CPU核心数配额';
COMMENT ON COLUMN resource_quotas.memory_quota IS '内存配额(MB)';
COMMENT ON COLUMN resource_quotas.gpu_quota IS 'GPU数量配额';
COMMENT ON COLUMN resource_quotas.storage_quota IS '存储配额(字节)';
COMMENT ON COLUMN resource_quotas.environment_quota IS '环境数量配额';
