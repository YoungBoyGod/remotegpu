# RemoteGPU 数据库 SQL 脚本

> 本目录包含 RemoteGPU 系统的所有数据库表结构 SQL 脚本
>
> 数据库类型：PostgreSQL 14+

---

## 📋 执行顺序

请按照以下顺序执行 SQL 脚本：

| 序号 | 文件名 | 说明 | 表数量 |
|------|--------|------|--------|
| 1 | `01_init_database.sql` | 数据库初始化、扩展、触发器函数 | 0 |
| 2 | `02_system_config.sql` | 系统配置表 | 1 |
| 3 | `03_users_and_permissions.sql` | 用户、工作空间、权限表 | 4 |
| 4 | `04_hosts_and_devices.sql` | 主机、GPU 设备表 | 2 |
| 5 | `05_environments.sql` | 开发环境、端口映射表 | 2 |
| 6 | `06_data_and_images.sql` | 数据集、模型、镜像表 | 6 |
| 7 | `07_monitoring.sql` | 监控数据表 | 3 |
| 8 | `08_billing.sql` | 计费、账单表 | 2 |
| 9 | `09_training_and_inference.sql` | 训练任务、推理服务表 | 2 |
| 10 | `10_notifications_and_logs.sql` | 通知、审计日志表 | 2 |
| 11 | `11_alerts_and_webhooks.sql` | 告警、Webhook 表 | 4 |
| 12 | `12_issues_and_requirements.sql` | 问题单、需求单、评论表 | 3 |
| 13 | `13_relationships.sql` | 关联关系、制品表 | 2 |

**总计：33 张表**

---

## 🚀 快速执行

### 方式一：逐个执行

```bash
psql -U postgres -d remotegpu -f 01_init_database.sql
psql -U postgres -d remotegpu -f 02_system_config.sql
psql -U postgres -d remotegpu -f 03_users_and_permissions.sql
# ... 依次执行其他文件
```

### 方式二：批量执行

```bash
#!/bin/bash
for i in {01..13}; do
    file=$(ls ${i}_*.sql 2>/dev/null)
    if [ -f "$file" ]; then
        echo "执行: $file"
        psql -U postgres -d remotegpu -f "$file"
    fi
done
```

### 方式三：合并执行

```bash
cat 0*.sql 1*.sql > all_tables.sql
psql -U postgres -d remotegpu -f all_tables.sql
```

---

## 📊 表结构概览

### 核心业务表

- **用户管理**: customers, workspaces, workspace_members, resource_quotas
- **设备管理**: hosts, gpus
- **环境管理**: environments, port_mappings
- **数据管理**: datasets, dataset_versions, models, model_versions
- **镜像管理**: images
- **训练推理**: training_jobs, inference_services

### 监控数据表

- **主机监控**: host_metrics
- **GPU监控**: gpu_metrics
- **环境监控**: environment_metrics

### 计费管理表

- **计费记录**: billing_records
- **账单**: invoices

### 辅助功能表

- **通知**: notifications
- **日志**: audit_logs
- **告警**: alert_rules, alert_records
- **Webhook**: webhooks, webhook_logs
- **工单**: issues, requirements, comments
- **制品**: artifacts
- **关联**: dataset_usage

---

## 🔧 设计特点

### 1. 减少外键依赖

为了提高灵活性和性能，本设计**尽量减少了外键约束**：

- 使用逻辑外键而非物理外键
- 通过应用层保证数据一致性
- 避免级联删除带来的性能问题

### 2. 索引优化

- 为常用查询字段创建索引
- 时序数据表使用复合索引（如 `host_id, collected_at DESC`）
- JSONB 字段使用 GIN 索引

### 3. 时间戳自动更新

- 所有表都有 `created_at` 字段
- 需要更新时间的表有 `updated_at` 字段和触发器

### 4. 软删除支持

- 部分表支持软删除（`deleted_at` 字段）
- 如：customers 表

---

## 📝 注意事项

### 1. 数据库创建

在执行 SQL 脚本前，需要先创建数据库：

```sql
CREATE DATABASE remotegpu
    WITH ENCODING 'UTF8'
    LC_COLLATE='en_US.UTF-8'
    LC_CTYPE='en_US.UTF-8';
```

### 2. 扩展依赖

需要安装以下 PostgreSQL 扩展：

- `uuid-ossp`: UUID 生成
- `pgcrypto`: 密码加密

### 3. 权限设置

建议创建专用数据库用户：

```sql
CREATE USER remotegpu_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE remotegpu TO remotegpu_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO remotegpu_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO remotegpu_user;
```

### 4. 时序数据保留策略

监控数据表（`*_metrics`）会产生大量数据，建议：

- 详细数据保留 7 天
- 聚合数据（小时级）保留 90 天
- 聚合数据（天级）保留 1 年

可以使用 PostgreSQL 的分区表或定期清理脚本。

---

## 🔍 表关系说明

### 核心关系

```
customers (客户)
  ├─> workspaces (工作空间)
  │     └─> workspace_members (成员)
  ├─> environments (环境)
  │     ├─> port_mappings (端口映射)
  │     └─> dataset_usage (数据集使用)
  ├─> datasets (数据集)
  ├─> models (模型)
  └─> billing_records (计费记录)

hosts (主机)
  ├─> gpus (GPU)
  ├─> host_metrics (主机监控)
  └─> environments (环境)
```

---

## 📖 相关文档

- [数据库设计文档](../docs/design/database_design.md)
- [客户管理设计](../docs/design/customer_management.md)
- [系统架构设计](../docs/design/system_architecture.md)

---

**创建日期**: 2026-01-26
**维护者**: RemoteGPU 团队
