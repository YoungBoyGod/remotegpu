-- ============================================
-- RemoteGPU 问题单和需求单表
-- ============================================
-- 文件: 12_issues_and_requirements.sql
-- 说明: 创建问题单、需求单相关表
-- 执行顺序: 12
-- ============================================

-- 问题单表
CREATE TABLE IF NOT EXISTS issues (
    id BIGSERIAL PRIMARY KEY,
    issue_no VARCHAR(64) UNIQUE NOT NULL,
    customer_id BIGINT NOT NULL,
    title VARCHAR(256) NOT NULL,
    description TEXT,
    type VARCHAR(32) NOT NULL,
    priority VARCHAR(20) DEFAULT 'medium',
    status VARCHAR(20) DEFAULT 'open',
    assignee_id BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    resolved_at TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_issues_customer ON issues(customer_id);
CREATE INDEX idx_issues_status ON issues(status);
CREATE INDEX idx_issues_assignee ON issues(assignee_id);
CREATE INDEX idx_issues_created_at ON issues(created_at DESC);

-- 创建更新时间触发器
CREATE TRIGGER update_issues_updated_at
    BEFORE UPDATE ON issues
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE issues IS '问题单表';
COMMENT ON COLUMN issues.type IS '类型: bug-缺陷, question-疑问, feature-功能请求';
COMMENT ON COLUMN issues.priority IS '优先级: low-低, medium-中, high-高, critical-紧急';
COMMENT ON COLUMN issues.status IS '状态: open-打开, in_progress-处理中, resolved-已解决, closed-已关闭';

-- 需求单表
CREATE TABLE IF NOT EXISTS requirements (
    id BIGSERIAL PRIMARY KEY,
    requirement_no VARCHAR(64) UNIQUE NOT NULL,
    customer_id BIGINT NOT NULL,
    title VARCHAR(256) NOT NULL,
    description TEXT,
    priority VARCHAR(20) DEFAULT 'medium',
    status VARCHAR(20) DEFAULT 'submitted',
    assignee_id BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_requirements_customer ON requirements(customer_id);
CREATE INDEX idx_requirements_status ON requirements(status);
CREATE INDEX idx_requirements_assignee ON requirements(assignee_id);
CREATE INDEX idx_requirements_created_at ON requirements(created_at DESC);

-- 创建更新时间触发器
CREATE TRIGGER update_requirements_updated_at
    BEFORE UPDATE ON requirements
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE requirements IS '需求单表';
COMMENT ON COLUMN requirements.priority IS '优先级: low-低, medium-中, high-高, critical-紧急';
COMMENT ON COLUMN requirements.status IS '状态: submitted-已提交, reviewing-审核中, approved-已批准, in_progress-开发中, completed-已完成, rejected-已拒绝';

-- 评论表（用于问题单和需求单）
CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    resource_type VARCHAR(32) NOT NULL,
    resource_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_comments_resource ON comments(resource_type, resource_id);
CREATE INDEX idx_comments_customer ON comments(customer_id);
CREATE INDEX idx_comments_created_at ON comments(created_at DESC);

-- 创建更新时间触发器
CREATE TRIGGER update_comments_updated_at
    BEFORE UPDATE ON comments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE comments IS '评论表';
COMMENT ON COLUMN comments.resource_type IS '资源类型: issue-问题单, requirement-需求单';
