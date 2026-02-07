# 数据库 Schema 审查和索引优化报告

> 审查人：architect
> 日期：2026-02-07
> 范围：backend/sql/ 下全部 28 个迁移脚本

---

## 一、迁移脚本清单

| 编号 | 文件名 | 内容 |
|------|--------|------|
| 01 | init_database.sql | 数据库初始化、扩展、触发器 |
| 02 | system_config.sql | 系统配置表 |
| 03 | users_and_permissions.sql | 客户、工作空间、配额 |
| 04 | hosts_and_devices.sql | 主机和 GPU 设备 |
| 05 | environments.sql | 开发环境和端口映射 |
| 06 | data_and_images.sql | 数据集、模型、镜像 |
| 07 | monitoring.sql | GPU 和环境监控 |
| 08 | billing.sql | 计费记录和账单 |
| 09 | training_and_inference.sql | 训练任务和推理服务 |
| 10 | notifications_and_logs.sql | 通知和审计日志 |
| 11 | alerts_and_webhooks.sql | 告警规则和 Webhook |
| 12 | issues_and_requirements.sql | 问题单和需求单 |
| 13 | relationships.sql | 数据集使用和制品 |
| 14 | machine_enrollments.sql | 机器添加任务 |
| 15 | task_process_id.sql | 任务进程 ID |
| 16 | add_customer_must_change_password.sql | 密码变更标志 |
| 17 | add_customer_company_code.sql | 公司代码 |
| 18 | add_host_metrics.sql | 主机监控指标 |
| 19 | add_ssh_info.sql | SSH 连接信息 |
| 20 | split_machine_status.sql | 拆分机器状态 |
| 21 | add_jupyter_vnc_info.sql | Jupyter/VNC 信息 |
| 22 | documents.sql | 文档中心 |
| 23 | add_audit_log_detail.sql / add_customer_quota.sql | 审计详情/客户配额 |
| 24 | add_task_output_progress.sql | 任务输出和进度 |
| 25 | add_system_config_group.sql | 配置分组 |
| 26 | dataset_mounts.sql | 数据集挂载 |
| 27 | seed_default_images.sql | 镜像种子数据 |
| 28 | add_host_external_mapping.sql | 外映射配置 |

---

## 二、索引覆盖情况总结

| 表名 | 已有索引 | 缺失索引 | 缺失复合索引 | 优先级 |
|------|---------|---------|------------|--------|
| hosts | 7 | 4 | 0 | 高 |
| allocations | 3 | 0 | 3 | 高 |
| tasks | 0 | 2 | 1 | 高 |
| dataset_mounts | 3 | 0 | 3 | 高 |
| notifications | 3 | 0 | 1 | 中 |
| customers | 5 | 1 | 0 | 中 |
| gpus | 4 | 0 | 0 | 低 |
| images | 4 | 0 | 0 | 低 |
| documents | 2 | 0 | 0 | 低 |

---

## 三、高优先级问题详解

### 3.1 hosts 表缺失索引

DAO 中频繁查询但缺少索引的字段：

| 字段 | 查询方法 | 说明 |
|------|---------|------|
| ip_address | FindByIPAddress() | 创建/导入时唯一性校验 |
| hostname | FindByHostname() | 创建/导入时唯一性校验 |
| region | List() 过滤 | 管理端列表筛选 |
| needs_collect | ListNeedCollect() | 采集任务查询 |
| last_heartbeat | HeartbeatMonitor | 心跳超时检测排序 |

### 3.2 tasks 表完全缺失索引

tasks 表是核心业务表，但 SQL 迁移中**没有定义任何索引**。

| 缺失字段 | 查询方法 | 说明 |
|---------|---------|------|
| customer_id | ListByCustomerID() | 客户查看自己的任务 |
| machine_id + status | ClaimTasks() | Agent 认领任务（原子操作） |
| status | ListAll() 过滤 | 管理端任务列表 |

### 3.3 allocations 表缺失复合索引

allocations 表的查询几乎都是多条件组合：

| 缺失复合索引 | 查询方法 |
|-------------|---------|
| (customer_id, status) | FindAllActiveByCustomerID() |
| (host_id, status) | FindActiveByHostID() |
| (host_id, customer_id, status) | FindActiveByHostAndCustomer() |

### 3.4 dataset_mounts 表缺失复合索引

| 缺失复合索引 | 查询方法 |
|-------------|---------|
| (dataset_id, status) | ListByDatasetID() |
| (host_id, status) | ListByHostID()、CountByHostID() |
| (dataset_id, host_id, status) | FindActiveMount() |

---

## 四、中优先级问题

### 4.1 notifications 表缺失复合索引

- 缺失 `(customer_id, is_read)` 复合索引
- 影响：ListByCustomerID()、CountUnread()、MarkAllRead() 三个高频查询

### 4.2 customers 表缺失 role 索引

- GORM 实体定义了 `index` 标签，但 SQL 迁移中未创建
- 影响：按角色筛选客户列表

---

## 五、外键约束问题

### 5.1 documents 表外键缺失 ON DELETE 策略

```sql
-- 当前定义
uploaded_by INT REFERENCES customers(id)
-- 缺少 ON DELETE SET NULL 或 ON DELETE CASCADE
```

如果删除客户，会导致外键约束违反。

### 5.2 SSH 字段重复定义

- `14_machine_enrollments.sql` 添加了 ssh_username/ssh_password/ssh_key
- `19_add_ssh_info.sql` 又添加了相同字段
- 可能导致迁移冲突，需要确认执行顺序

---

## 六、建议的优化迁移脚本

建议创建 `29_add_missing_indexes.sql` 统一添加缺失索引。
