# 数据库设计变更摘要

> 本文档记录Phase 0.5数据库设计修正的所有变更
>
> **修正日期**: 2026-01-28
> **修正原因**: 支持平台核心需求（多挂载、工作空间计费、扩缩容等）
> **版本**: 1.0

---

## 📋 变更概览

| SQL文件 | 修改的表 | 新增字段数 | 删除字段数 | 影响 |
|---------|---------|-----------|-----------|------|
| 05_environments.sql | environments | 15 | 0 | 🔴 高 |
| 05_environments.sql | port_mappings | 3 | 0 | 🟡 中 |
| 09_training_and_inference.sql | training_jobs | 12 | 3 | 🔴 高 |
| 09_training_and_inference.sql | inference_services | 13 | 0 | 🔴 高 |
| 08_billing.sql | billing_records | 6 | 1 | 🔴 高 |
| 03_users_and_permissions.sql | customers | 5 | 0 | 🔴 高 |

**总计:** 6个表，新增54个字段，删除4个字段

---

## 1. environments表变更

### 1.1 新增字段（15个）

**资源配置（2个）:**
- `gpu_memory BIGINT` - GPU显存(字节)
- `temp_storage BIGINT` - 临时存储空间(字节)

**访问配置（7个）:**
- `ssh_enabled BOOLEAN DEFAULT true` - SSH访问开关
- `rdp_enabled BOOLEAN DEFAULT false` - RDP访问开关
- `jupyter_token VARCHAR(128)` - JupyterLab访问令牌
- `jupyter_enabled BOOLEAN DEFAULT true` - JupyterLab开关
- `tensorboard_port INT` - TensorBoard端口
- `tensorboard_enabled BOOLEAN DEFAULT false` - TensorBoard开关
- `web_terminal_enabled BOOLEAN DEFAULT true` - Web Terminal开关

**挂载配置（2个）:**
- `mounted_datasets JSONB` - 挂载的数据集列表 `[{"id": 1, "path": "/gemini/data-1", "readonly": true}]`
- `mounted_models JSONB` - 挂载的模型列表 `[{"id": 1, "path": "/gemini/pretrain1", "readonly": true}]`

**环境配置（2个）:**
- `env_vars JSONB` - 环境变量 `{"KEY": "value"}`
- `config JSONB` - 其他配置

**软删除（1个）:**
- `deleted_at TIMESTAMP` - 软删除时间

### 1.2 新增索引（5个）
- `idx_environments_deleted_at` - 软删除索引
- `idx_environments_mounted_datasets` (GIN) - JSONB索引
- `idx_environments_mounted_models` (GIN) - JSONB索引
- `idx_environments_env_vars` (GIN) - JSONB索引
- `idx_environments_config` (GIN) - JSONB索引

---

## 2. port_mappings表变更

### 2.1 新增字段（3个）

**生命周期管理:**
- `last_accessed_at TIMESTAMP` - 最后访问时间
- `auto_release_hours INT DEFAULT 48` - 自动释放时间(小时)

### 2.2 更新注释
- `service_type` 增加 `tensorboard` 类型

---

## 3. training_jobs表变更

### 3.1 删除字段（3个）
- ❌ `dataset_id BIGINT` - 改为JSONB数组支持多挂载
- ❌ `model_id BIGINT` - 改为JSONB数组支持多挂载
- ❌ `script_path VARCHAR(512)` - 改为command TEXT

### 3.2 新增字段（12个）

**镜像和框架:**
- `image VARCHAR(256) NOT NULL` - 镜像名称

**挂载配置（2个）:**
- `mounted_datasets JSONB` - 挂载的数据集ID列表 `[1, 2, 3]` 最多3个
- `mounted_models JSONB` - 挂载的模型ID列表 `[1, 2]` 最多3个

**执行配置（3个）:**
- `command TEXT NOT NULL` - 启动命令（替代script_path）
- `env_vars JSONB` - 环境变量
- `output_path VARCHAR(512) DEFAULT '/gemini/output/'` - 输出路径

**分布式训练（2个）:**
- `node_count INT DEFAULT 1` - 分布式训练节点数量
- `distributed_config JSONB` - 分布式训练配置 `{"framework": "pytorch", "backend": "nccl"}`

**资源配置（1个）:**
- `gpu_memory BIGINT` - GPU显存(字节)

**软删除（1个）:**
- `deleted_at TIMESTAMP` - 软删除时间

### 3.3 新增索引（5个）
- `idx_training_jobs_deleted_at` - 软删除索引
- `idx_training_jobs_mounted_datasets` (GIN) - JSONB索引
- `idx_training_jobs_mounted_models` (GIN) - JSONB索引
- `idx_training_jobs_env_vars` (GIN) - JSONB索引
- `idx_training_jobs_distributed_config` (GIN) - JSONB索引

### 3.4 更新注释
- `status` 增加 `queued-排队中` 状态

---

## 4. inference_services表变更

### 4.1 新增字段（13个）

**镜像:**
- `image VARCHAR(256) NOT NULL` - 镜像名称

**副本配置（2个）:**
- `min_replicas INT DEFAULT 1` - 最小副本数量
- `max_replicas INT DEFAULT 10` - 最大副本数量

**自动扩缩容（2个）:**
- `autoscaling_enabled BOOLEAN DEFAULT false` - 是否启用自动扩缩容
- `autoscaling_config JSONB` - 扩缩容配置 `{"target_cpu": 80, "target_qps": 1000}`

**版本管理（2个）:**
- `version VARCHAR(64)` - 当前版本
- `previous_version VARCHAR(64)` - 上一个版本（用于回滚）

**健康检查（2个）:**
- `health_check_path VARCHAR(256)` - 健康检查路径
- `health_check_interval INT DEFAULT 30` - 健康检查间隔(秒)

**环境配置（1个）:**
- `env_vars JSONB` - 环境变量

**资源配置（1个）:**
- `gpu_memory BIGINT` - GPU显存(字节)

**软删除（1个）:**
- `deleted_at TIMESTAMP` - 软删除时间

### 4.2 新增索引（3个）
- `idx_inference_services_deleted_at` - 软删除索引
- `idx_inference_services_autoscaling_config` (GIN) - JSONB索引
- `idx_inference_services_env_vars` (GIN) - JSONB索引

---

## 5. billing_records表变更

### 5.1 删除字段（1个）
- ❌ `env_id VARCHAR(64)` - 改为resource_id（更通用）

### 5.2 新增字段（6个）

**工作空间维度:**
- `workspace_id BIGINT` - 工作空间ID（工作空间维度计费）

**资源关联（2个）:**
- `resource_id VARCHAR(64)` - 资源ID（env_id, training_job_id, inference_service_id等）
- `resource_name VARCHAR(256)` - 资源名称

**计费粒度（2个）:**
- `billing_unit VARCHAR(20) DEFAULT 'minute'` - 计费单位: minute, hour, day
- `billing_minutes INT` - 计费分钟数（向上取整）

### 5.3 更新索引
- ❌ 删除 `idx_billing_records_env`
- ✅ 新增 `idx_billing_records_workspace`
- ✅ 新增 `idx_billing_records_resource`

### 5.4 更新注释
- `resource_type` 改为: environment, training, inference, storage

---

## 6. customers表变更

### 6.1 新增字段（5个）

**账户余额（2个）:**
- `balance DECIMAL(10,4) DEFAULT 0.00` - 账户余额（算力点）
- `currency VARCHAR(10) DEFAULT 'CNY'` - 货币类型

**欠费管理（3个）:**
- `overdue_status VARCHAR(20) DEFAULT 'normal'` - 欠费状态: normal, overdue, suspended
- `overdue_since TIMESTAMP` - 欠费开始时间
- `last_payment_at TIMESTAMP` - 最后充值时间

### 6.2 新增索引（1个）
- `idx_customers_overdue_status` - 欠费状态索引

---

## 7. 兼容性说明

### 7.1 破坏性变更

**training_jobs表:**
- ❌ 删除 `dataset_id` 字段 → 使用 `mounted_datasets` JSONB数组
- ❌ 删除 `model_id` 字段 → 使用 `mounted_models` JSONB数组
- ❌ 删除 `script_path` 字段 → 使用 `command` TEXT字段

**billing_records表:**
- ❌ 删除 `env_id` 字段 → 使用 `resource_id` 字段

### 7.2 迁移建议

如果已有数据，需要执行数据迁移：

```sql
-- training_jobs表迁移
UPDATE training_jobs
SET mounted_datasets = jsonb_build_array(dataset_id)
WHERE dataset_id IS NOT NULL;

UPDATE training_jobs
SET mounted_models = jsonb_build_array(model_id)
WHERE model_id IS NOT NULL;

UPDATE training_jobs
SET command = script_path
WHERE script_path IS NOT NULL;

-- billing_records表迁移
UPDATE billing_records
SET resource_id = env_id
WHERE env_id IS NOT NULL;
```

---

## 8. 验证清单

### 8.1 SQL文件验证
- [x] 05_environments.sql - 语法正确
- [x] 09_training_and_inference.sql - 语法正确
- [x] 08_billing.sql - 语法正确
- [x] 03_users_and_permissions.sql - 语法正确

### 8.2 功能验证
- [ ] environments表支持4种访问方式配置
- [ ] training_jobs表支持3个数据集+3个模型挂载
- [ ] billing_records表支持工作空间维度计费
- [ ] customers表支持余额和欠费管理
- [ ] inference_services表支持扩缩容配置

### 8.3 索引验证
- [x] 所有JSONB字段已创建GIN索引
- [x] 所有deleted_at字段已创建索引
- [x] 所有外键字段已创建索引

---

## 9. 下一步行动

### 9.1 立即行动
1. ✅ 修正SQL文件 - 已完成
2. ⏳ 更新数据库设计文档
3. ⏳ 实现Go实体（Phase 1A）

### 9.2 后续工作
4. 测试数据库迁移脚本
5. 验证表结构完整性
6. 更新API文档

---

## 10. 影响评估

### 10.1 正面影响
- ✅ 支持核心需求（多挂载、工作空间计费）
- ✅ 支持高级功能（扩缩容、版本回滚）
- ✅ 提升数据查询性能（GIN索引）
- ✅ 支持软删除（数据安全）

### 10.2 注意事项
- ⚠️ JSONB字段需要应用层验证数据格式
- ⚠️ 破坏性变更需要数据迁移
- ⚠️ GIN索引会增加写入开销

---

**文档维护者:** RemoteGPU开发团队
**最后更新:** 2026-01-28
