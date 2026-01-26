# 客户管理设计文档

> 本文档描述 RemoteGPU 系统的客户管理模块设计
>
> 创建日期：2026-01-26

---

## 1. 概述

### 1.1 模块定位

客户管理模块是 RemoteGPU 系统的核心基础模块之一，负责管理用户账户、认证、授权、工作空间和配额等功能。

### 1.2 核心功能

- 用户注册与登录
- 用户信息管理
- 工作空间管理
- 资源配额管理
- 权限控制（RBAC）

---

## 2. 数据模型设计

### 2.1 用户表（customers）

```sql
CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    username VARCHAR(64) UNIQUE NOT NULL,
    email VARCHAR(128) UNIQUE NOT NULL,
    password_hash VARCHAR(256) NOT NULL,
    display_name VARCHAR(128),
    avatar_url VARCHAR(512),
    phone VARCHAR(32),
    account_type VARCHAR(32) DEFAULT 'individual',
    status VARCHAR(32) DEFAULT 'active',
    email_verified BOOLEAN DEFAULT FALSE,
    phone_verified BOOLEAN DEFAULT FALSE,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_customers_username ON customers(username);
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_status ON customers(status);
```


**字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGSERIAL | 主键 |
| uuid | UUID | 全局唯一标识符 |
| username | VARCHAR(64) | 用户名（唯一） |
| email | VARCHAR(128) | 邮箱（唯一） |
| password_hash | VARCHAR(256) | 密码哈希（bcrypt） |
| display_name | VARCHAR(128) | 显示名称 |
| avatar_url | VARCHAR(512) | 头像 URL |
| phone | VARCHAR(32) | 手机号 |
| account_type | VARCHAR(32) | 账户类型（individual/enterprise） |
| status | VARCHAR(32) | 状态（active/suspended/deleted） |
| email_verified | BOOLEAN | 邮箱是否已验证 |
| phone_verified | BOOLEAN | 手机号是否已验证 |
| last_login_at | TIMESTAMP | 最后登录时间 |


### 2.2 工作空间表（workspaces）

```sql
CREATE TABLE workspaces (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    owner_id BIGINT NOT NULL,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    type VARCHAR(32) DEFAULT 'personal',
    member_count INT DEFAULT 1,
    status VARCHAR(32) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (owner_id) REFERENCES customers(id)
);

CREATE INDEX idx_workspaces_owner ON workspaces(owner_id);
CREATE INDEX idx_workspaces_type ON workspaces(type);
```

**字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| type | VARCHAR(32) | 类型（personal/team/enterprise） |
| member_count | INT | 成员数量 |
| status | VARCHAR(32) | 状态（active/archived） |


### 2.3 工作空间成员表（workspace_members）

```sql
CREATE TABLE workspace_members (
    id BIGSERIAL PRIMARY KEY,
    workspace_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    role VARCHAR(32) DEFAULT 'member',
    status VARCHAR(32) DEFAULT 'active',
    joined_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    UNIQUE(workspace_id, customer_id)
);
```

**角色定义：**
- `owner` - 所有者（完全控制权）
- `admin` - 管理员（管理成员和资源）
- `member` - 成员（使用资源）
- `viewer` - 查看者（只读权限）


### 2.4 资源配额表（resource_quotas）

```sql
CREATE TABLE resource_quotas (
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
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id)
);
```


**配额级别：**

| 级别 | CPU | 内存 | GPU | 存储 | 环境数 |
|------|-----|------|-----|------|--------|
| Free | 4 核 | 8GB | 0 | 10GB | 1 |
| Basic | 8 核 | 16GB | 1 | 50GB | 3 |
| Pro | 16 核 | 32GB | 2 | 200GB | 10 |
| Enterprise | 自定义 | 自定义 | 自定义 | 自定义 | 无限 |

---

## 3. 用户认证设计

### 3.1 认证方式

**支持的认证方式：**
1. 用户名密码认证
2. 邮箱密码认证
3. OAuth2 认证（GitHub, Google, GitLab）
4. LDAP/AD 认证（企业版）


### 3.2 JWT Token 设计

**Token 结构：**
```json
{
  "customer_id": 123,
  "username": "user@example.com",
  "workspace_id": 456,
  "role": "member",
  "exp": 1706256000,
  "iat": 1706169600
}
```

**Token 有效期：**
- Access Token: 24 小时
- Refresh Token: 30 天


### 3.3 密码安全策略

**密码要求：**
- 最小长度：8 位
- 必须包含：大写字母、小写字母、数字
- 建议包含：特殊字符
- 密码哈希：使用 bcrypt（cost=10）

**安全措施：**
- 登录失败 5 次后锁定账户 15 分钟
- 密码定期更换提醒（90 天）
- 支持双因素认证（2FA）

---

## 4. 权限控制（RBAC）

### 4.1 权限模型

```
用户 (Customer)
  └─ 工作空间成员 (Workspace Member)
      └─ 角色 (Role)
          └─ 权限 (Permission)
```


### 4.2 权限列表

**环境管理权限：**
- `environment:create` - 创建环境
- `environment:read` - 查看环境
- `environment:update` - 修改环境
- `environment:delete` - 删除环境
- `environment:start` - 启动环境
- `environment:stop` - 停止环境

**数据集权限：**
- `dataset:create` - 创建数据集
- `dataset:read` - 查看数据集
- `dataset:update` - 修改数据集
- `dataset:delete` - 删除数据集

**工作空间权限：**
- `workspace:manage_members` - 管理成员
- `workspace:manage_settings` - 管理设置


### 4.3 角色权限矩阵

| 权限 | Owner | Admin | Member | Viewer |
|------|-------|-------|--------|--------|
| 创建环境 | ✅ | ✅ | ✅ | ❌ |
| 查看环境 | ✅ | ✅ | ✅ | ✅ |
| 修改环境 | ✅ | ✅ | ✅ | ❌ |
| 删除环境 | ✅ | ✅ | ✅ | ❌ |
| 管理成员 | ✅ | ✅ | ❌ | ❌ |
| 管理设置 | ✅ | ❌ | ❌ | ❌ |

---

## 5. API 接口设计

### 5.1 用户注册

```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "user123",
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```


**响应：**
```json
{
  "customer_id": 123,
  "username": "user123",
  "email": "user@example.com",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### 5.2 用户登录

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "user123",
  "password": "SecurePass123!"
}
```


### 5.3 获取用户信息

```http
GET /api/users/me
Authorization: Bearer {token}
```

**响应：**
```json
{
  "customer_id": 123,
  "username": "user123",
  "email": "user@example.com",
  "display_name": "User Name",
  "avatar_url": "https://...",
  "account_type": "individual",
  "quota": {
    "cpu": 8,
    "memory": 16384,
    "gpu": 1,
    "storage": 53687091200
  }
}
```


### 5.4 创建工作空间

```http
POST /api/workspaces
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "My Team",
  "description": "Team workspace",
  "type": "team"
}
```

### 5.5 添加工作空间成员

```http
POST /api/workspaces/{workspace_id}/members
Authorization: Bearer {token}
Content-Type: application/json

{
  "customer_id": 456,
  "role": "member"
}
```


---

## 6. 业务流程

### 6.1 用户注册流程

```
1. 用户提交注册信息
2. 验证用户名和邮箱唯一性
3. 密码强度检查
4. 创建用户记录（密码 bcrypt 加密）
5. 创建默认个人工作空间
6. 分配默认配额（Free 级别）
7. 发送邮箱验证邮件
8. 返回 JWT Token
```

### 6.2 工作空间创建流程

```
1. 验证用户权限
2. 检查工作空间名称唯一性
3. 创建工作空间记录
4. 添加创建者为 Owner
5. 分配工作空间配额
6. 返回工作空间信息
```


---

## 7. 安全考虑

### 7.1 数据安全

- 密码使用 bcrypt 加密存储
- 敏感信息（如手机号）加密存储
- API 接口使用 HTTPS
- JWT Token 签名验证

### 7.2 访问控制

- 所有 API 需要认证
- 基于 RBAC 的权限控制
- 工作空间资源隔离
- 操作审计日志


### 7.3 用户分类

**用户类型定义：**

1. **管理人员（Admin）**
   - 系统管理员，拥有最高权限
   - 可以管理所有用户、工作空间和资源
   - 可以配置系统设置和基础设施
   - 可以查看所有审计日志

2. **内部用户（Internal User）**
   - 企业内部员工
   - 通过 LDAP/AD 或内部邮箱域名识别
   - 拥有较高的资源配额
   - 可以创建团队工作空间

3. **外部用户（External User）**
   - 公开注册的普通用户
   - 通过邮箱注册
   - 默认较低的资源配额
   - 只能使用个人工作空间


**用户表更新（添加 user_type 字段）：**

```sql
ALTER TABLE customers ADD COLUMN user_type VARCHAR(32) DEFAULT 'external';
CREATE INDEX idx_customers_user_type ON customers(user_type);
```

**用户类型枚举值：**
- `admin` - 管理人员
- `internal` - 内部用户
- `external` - 外部用户


**不同用户类型的权限差异：**

| 功能 | 管理人员 | 内部用户 | 外部用户 |
|------|---------|---------|---------|
| 管理所有用户 | ✅ | ❌ | ❌ |
| 管理系统设置 | ✅ | ❌ | ❌ |
| 查看审计日志 | ✅ | ❌ | ❌ |
| 创建团队工作空间 | ✅ | ✅ | ❌ |
| 创建个人工作空间 | ✅ | ✅ | ✅ |
| 使用 GPU 资源 | ✅ | ✅ | ✅ |


**不同用户类型的默认配额：**

| 资源 | 管理人员 | 内部用户 | 外部用户 |
|------|---------|---------|---------|
| CPU 核心 | 无限 | 16 核 | 4 核 |
| 内存 | 无限 | 32GB | 8GB |
| GPU 卡数 | 无限 | 2 张 | 0 张 |
| 存储空间 | 无限 | 200GB | 10GB |
| 环境数量 | 无限 | 10 个 | 1 个 |


**用户类型识别规则：**

1. **管理人员识别：**
   - 手动在数据库中设置 `user_type = 'admin'`
   - 或通过管理后台界面设置

2. **内部用户识别：**
   - 邮箱域名在白名单中（如 @company.com）
   - 通过 LDAP/AD 认证的用户
   - 注册时自动识别并设置

3. **外部用户识别：**
   - 默认类型
   - 公开邮箱注册的用户


---

## 8. 实现建议

### 8.1 用户注册时的类型判断

```go
func DetermineUserType(email string) string {
    // 检查是否为内部邮箱域名
    internalDomains := []string{"@company.com", "@internal.com"}
    for _, domain := range internalDomains {
        if strings.HasSuffix(email, domain) {
            return "internal"
        }
    }
    return "external"
}
```


### 8.2 权限检查中间件

```go
func CheckUserType(requiredType string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userType := c.GetString("user_type")
        
        if requiredType == "admin" && userType != "admin" {
            c.JSON(403, gin.H{"error": "需要管理员权限"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```


### 8.3 配额分配策略

```go
func AssignDefaultQuota(userType string) ResourceQuota {
    quotaMap := map[string]ResourceQuota{
        "admin": {
            CPU:         -1,  // 无限
            Memory:      -1,
            GPU:         -1,
            Storage:     -1,
            Environment: -1,
        },
        "internal": {
            CPU:         16,
            Memory:      32768,
            GPU:         2,
            Storage:     214748364800,  // 200GB
            Environment: 10,
        },
        "external": {
            CPU:         4,
            Memory:      8192,
            GPU:         0,
            Storage:     10737418240,  // 10GB
            Environment: 1,
        },
    }
    return quotaMap[userType]
}
```


---

## 9. 总结

本文档定义了 RemoteGPU 系统的客户管理模块设计，包括：

1. **用户分类体系**
   - 管理人员：系统管理员，拥有最高权限
   - 内部用户：企业内部员工，较高配额
   - 外部用户：公开注册用户，基础配额

2. **核心功能**
   - 用户注册与认证
   - 工作空间管理
   - 资源配额管理
   - 基于 RBAC 的权限控制

3. **安全机制**
   - 密码 bcrypt 加密
   - JWT Token 认证
   - 权限分级控制
   - 操作审计日志

---

**文档版本：** v1.0  
**创建日期：** 2026-01-26  
**维护者：** RemoteGPU 团队

