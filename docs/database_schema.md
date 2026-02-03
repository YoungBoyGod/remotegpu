# RemoteGPU 数据库 Schema 详解 (V2.1)

本文档基于 `sql/remotegpu_complete_schema.sql` 生成，详细解释了每个表的设计意图、字段含义以及数据插入示例。

---

## 1. 系统配置 (system_configs)

用于存储全局系统设置，如“默认配额”、“站点名称”等。

| 字段 | 类型 | 说明 |
| :--- | :--- | :--- |
| `id` | BIGSERIAL | 主键 |
| `config_key` | VARCHAR(128) | 配置键（唯一），如 `site_name` |
| `config_value` | TEXT | 配置值 |
| `config_type` | VARCHAR(32) | 值类型：string, int, bool, json |
| `description` | TEXT | 配置说明 |
| `is_public` | BOOLEAN | 是否对前端公开 |

**插入示例**:
```sql
INSERT INTO system_configs (config_key, config_value, config_type, description, is_public) 
VALUES 
('site_name', 'RemoteGPU Cloud', 'string', '平台名称', true),
('default_gpu_quota', '2', 'int', '默认GPU配额', false);
```

---

## 2. 客户与租户 (customers)

核心用户表。在 V2.0 中，`Customer` 既代表“登录账号”，也代表“租户/公司”。

| 字段 | 类型 | 说明 |
| :--- | :--- | :--- |
| `id` | BIGSERIAL | 内部主键（用于外键关联） |
| `uuid` | UUID | 对外暴露的唯一 ID |
| `username` | VARCHAR(64) | 登录用户名 |
| `email` | VARCHAR(128) | 联系邮箱 |
| `password_hash` | VARCHAR(256) | 加密密码 |
| `role` | VARCHAR(32) | 角色：`admin` (超管), `customer_owner` (客户管理员), `customer_member` (普通成员) |
| `user_type` | VARCHAR(32) | 类型：`admin` (内部人员), `external` (外部客户) |
| `balance` | DECIMAL | 账户余额 |

**插入示例**:
```sql
-- 插入一个管理员
INSERT INTO customers (username, email, password_hash, role, user_type, display_name)
VALUES ('admin', 'admin@remotegpu.com', '$2a$10$xyz...', 'admin', 'admin', 'System Admin');

-- 插入一个企业客户
INSERT INTO customers (username, email, password_hash, role, user_type, company, balance)
VALUES ('acme_corp', 'it@acme.com', '$2a$10$abc...', 'customer_owner', 'external', 'Acme Inc.', 1000.00);
```

---

## 3. SSH 密钥 (ssh_keys)

客户上传的 SSH 公钥，用于自动注入到分配的机器中。

| 字段 | 类型 | 说明 |
| :--- | :--- | :--- |
| `customer_id` | BIGINT | 关联客户 ID |
| `name` | VARCHAR(64) | 密钥名称（如 "Office Laptop"） |
| `public_key` | TEXT | SSH 公钥内容 (ssh-rsa AAAA...) |
| `fingerprint` | VARCHAR(128)| 密钥指纹 |

**插入示例**:
```sql
INSERT INTO ssh_keys (customer_id, name, public_key, fingerprint)
VALUES (2, 'My Macbook', 'ssh-rsa AAAA...', 'SHA256:...');
```

---

## 4. 主机资源 (hosts)

物理机或虚拟机节点信息。

| 字段 | 类型 | 说明 |
| :--- | :--- | :--- |
| `id` | VARCHAR(64) | 机器 ID (如 "node-01")，手动指定或自动生成 |
| `ip_address` | VARCHAR(64) | 内网 IP (Agent 通信) |
| `public_ip` | VARCHAR(64) | 公网 IP (用户连接用) |
| `status` | VARCHAR(20) | `online`, `offline`, `allocated`, `maintenance` |
| `total_gpu` | INT | GPU 数量 |
| `region` | VARCHAR(64) | 地区 (如 "beijing-zone-a") |

**插入示例**:
```sql
INSERT INTO hosts (id, name, ip_address, public_ip, total_cpu, total_memory_gb, total_disk_gb, status, region)
VALUES ('node-01', 'BJ-GPU-01', '192.168.1.100', '203.0.113.10', 64, 256, 2048, 'online', 'beijing');
```

---

## 5. GPU 设备 (gpus)

主机上的具体 GPU 卡信息。

| 字段 | 类型 | 说明 |
| :--- | :--- | :--- |
| `host_id` | VARCHAR(64) | 关联主机 |
| `index` | INT | 显卡序号 (0, 1, 2...) |
| `name` | VARCHAR(128) | 型号 (如 "RTX 4090") |
| `memory_total_mb`| INT | 显存大小 (24576) |
| `status` | VARCHAR(20) | `available`, `allocated` |
| `allocated_to` | VARCHAR(64) | 占用者信息 (可选) |

**插入示例**:
```sql
INSERT INTO gpus (host_id, index, name, memory_total_mb, status)
VALUES ('node-01', 0, 'NVIDIA GeForce RTX 4090', 24576, 'available');
```

---

## 6. 资源分配 (allocations)

**核心表**：记录“谁”租了“哪台机器”。

| 字段 | 类型 | 说明 |
| :--- | :--- | :--- |
| `id` | VARCHAR(64) | 分配单号 (如 "alloc-20260203-001") |
| `customer_id` | BIGINT | 客户 ID |
| `host_id` | VARCHAR(64) | 机器 ID |
| `start_time` | TIMESTAMP | 开始时间 |
| `end_time` | TIMESTAMP | 预计结束时间 |
| `status` | VARCHAR(32) | `active` (使用中), `expired` (已到期), `reclaimed` (已回收) |

**插入示例**:
```sql
INSERT INTO allocations (id, customer_id, host_id, start_time, end_time, status, remark)
VALUES ('alloc-001', 2, 'node-01', NOW(), NOW() + INTERVAL '1 month', 'active', 'VIP客户测试');
```

---

## 7. 镜像 (images)

Docker 镜像元数据。

| 字段 | 类型 | 说明 |
| :--- | :--- | :--- |
| `name` | VARCHAR(256) | 镜像名 (如 "pytorch/2.0:latest") |
| `framework` | VARCHAR(64) | 框架 (PyTorch, TensorFlow) |
| `is_official` | BOOLEAN | 是否官方镜像 |
| `registry_url` | VARCHAR(512)| 仓库地址 |

**插入示例**:
```sql
INSERT INTO images (name, display_name, framework, cuda_version, is_official)
VALUES ('pytorch/2.0-cuda11.8', 'PyTorch 2.0 (CUDA 11.8)', 'PyTorch', '11.8', true);
```

---

## 8. 数据集 (datasets) & 挂载 (dataset_mounts)

**datasets**:
| 字段 | 类型 | 说明 |
| :--- | :--- | :--- |
| `storage_path` | VARCHAR(512)| MinIO/S3 上的桶路径 (如 "datasets/imagenet") |
| `status` | VARCHAR(20) | `uploading`, `ready` |

**dataset_mounts**:
| 字段 | 类型 | 说明 |
| :--- | :--- | :--- |
| `dataset_id` | BIGINT | 数据集 ID |
| `host_id` | VARCHAR(64) | 挂载到哪台机器 |
| `mount_path` | VARCHAR(256)| 机器内路径 (如 "/mnt/data/imagenet") |

**插入示例**:
```sql
-- 创建数据集记录
INSERT INTO datasets (customer_id, name, storage_path, status)
VALUES (2, 'ImageNet-Train', 's3://remotegpu/users/2/datasets/imagenet', 'ready');

-- 挂载记录
INSERT INTO dataset_mounts (dataset_id, host_id, mount_path)
VALUES (1, 'node-01', '/data/imagenet');
```

---

## 9. 任务 (tasks)

训练或推理任务。

| 字段 | 类型 | 说明 |
| :--- | :--- | :--- |
| `type` | VARCHAR(32) | `training`, `inference` |
| `command` | TEXT | 启动命令 (如 "python train.py") |
| `image_id` | BIGINT | 使用的镜像 |
| `status` | VARCHAR(20) | `queued`, `running`, `stopped` |

**插入示例**:
```sql
INSERT INTO tasks (id, customer_id, host_id, name, type, image_id, command, status)
VALUES ('task-101', 2, 'node-01', 'ResNet Training', 'training', 1, 'python train.py --epochs 100', 'queued');
```

---

## 10. 审计日志 (audit_logs)

记录敏感操作。

| 字段 | 类型 | 说明 |
| :--- | :--- | :--- |
| `action` | VARCHAR(128) | 操作动作 (如 "allocation.create") |
| `resource_type`| VARCHAR(64) | 资源类型 (如 "host") |
| `detail` | JSONB | 操作详情快照 |

**插入示例**:
```sql
INSERT INTO audit_logs (customer_id, username, action, resource_type, resource_id, detail)
VALUES (1, 'admin', 'allocation.create', 'host', 'node-01', '{"duration": "1m", "user": "acme"}');
```
