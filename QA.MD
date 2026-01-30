# QA 问答记录

本文档记录开发过程中的问题和解决方案。

---

## 2026-01-30: 存储模块上传功能测试

### 问题
需要为存储模块编写测试用例，测试上传功能到本地存储和 RustFS (S3兼容存储)。

### 解决方案

1. **创建测试文件**: `backend/pkg/storage/storage_test.go`

2. **本地存储测试** (`TestLocalStorage_Upload`):
   - 创建临时目录作为存储路径
   - 上传测试文件
   - 验证文件存在性和内容
   - 测试完成后自动清理

3. **S3存储测试** (`TestS3Storage_Upload`):
   - 支持通过环境变量配置凭证
   - 无配置时使用默认值，服务不可用时自动跳过
   - 测试上传、验证存在、下载验证内容
   - 测试完成后清理上传的文件

### 运行测试

```bash
# 本地存储测试
go test -v ./pkg/storage/... -run TestLocalStorage_Upload

# S3存储测试（需配置环境变量）
S3_ENDPOINT="http://192.168.10.210:9000" \
S3_ACCESS_KEY="EqyHX3vdSlfKj8hsiazR" \
S3_SECRET_KEY="19IEWuomU63LpKsVtHkJgPCcTx0evdfAOGDqwX7h" \
S3_BUCKET="remotegpu" \
go test -v ./pkg/storage/... -run TestS3Storage_Upload
```

### 相关文件
- `backend/pkg/storage/storage_test.go` - 测试文件
- `backend/pkg/storage/local.go` - 本地存储实现
- `backend/pkg/storage/s3.go` - S3存储实现

---

## 2026-01-30: 创建 hosts 和 gpus 数据库表

### 问题
开始开发主机管理模块，需要先创建数据库表结构。

### 解决方案

1. **执行初始化 SQL**：创建 `update_updated_at_column()` 触发器函数

2. **创建 hosts 表**：
   - 29 个字段，包含主机基本信息、资源配置、状态等
   - 5 个索引（os_type, status, deployment_mode, labels, health_status）
   - 自动更新 updated_at 的触发器

3. **创建 gpus 表**：
   - 18 个字段，包含 GPU 设备信息、状态、分配情况等
   - 4 个索引 + 2 个唯一约束（uuid, host_id+gpu_index）
   - 自动更新 updated_at 的触发器

### 执行命令

```bash
# 通过 Docker 执行 SQL
docker exec remotegpu-postgresql psql -U remotegpu_user -d remotegpu -c "..."

# 验证表结构
docker exec remotegpu-postgresql psql -U remotegpu_user -d remotegpu -c "\d hosts"
docker exec remotegpu-postgresql psql -U remotegpu_user -d remotegpu -c "\d gpus"
```

### 相关文件
- `backend/sql/01_init_database.sql` - 初始化脚本
- `backend/sql/04_hosts_and_devices.sql` - 主机和设备表定义

### 下一步
创建 Host 和 GPU 的 Entity、DAO、Service、Controller

---

## 2026-01-30: 实现 Host 和 GPU 完整功能

### 问题
需要实现主机和GPU管理的完整后端功能。

### 解决方案

**Entity 层**:
- `entity/host.go` - Host 实体 (29字段)
- `entity/gpu.go` - GPU 实体 (18字段)

**DAO 层**:
- `dao/host.go` - CRUD + 状态更新 + 心跳
- `dao/gpu.go` - CRUD + 按主机查询

**Service 层**:
- `service/host.go` - 业务逻辑封装
- `service/gpu.go` - 业务逻辑封装

**Controller 层**:
- `controller/v1/host.go` - REST API
- `controller/v1/gpu.go` - REST API

### API 路由

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /admin/hosts | 创建主机 |
| GET | /admin/hosts | 主机列表 |
| GET | /admin/hosts/:id | 主机详情 |
| PUT | /admin/hosts/:id | 更新主机 |
| DELETE | /admin/hosts/:id | 删除主机 |
| POST | /admin/hosts/:id/heartbeat | 心跳 |
| POST | /admin/gpus | 创建GPU |
| GET | /admin/gpus/:id | GPU详情 |
| GET | /admin/hosts/:host_id/gpus | 主机GPU列表 |
| DELETE | /admin/gpus/:id | 删除GPU |

---

## 2026-01-30: Host 和 GPU 模块测试完成

### 问题
需要为 Host 和 GPU 模块的各层编写测试用例。

### 解决方案

**DAO 层测试**:
- `dao/host_test.go` - Host DAO 测试 (CRUD)
- `dao/gpu_test.go` - GPU DAO 测试 (CRUD + 按主机查询)

**Service 层测试**:
- `service/host_test.go` - Host Service 测试
- `service/gpu_test.go` - GPU Service 测试

**Controller 层测试**:
- `controller/v1/host_test.go` - Host Controller HTTP 测试
- `controller/v1/gpu_test.go` - GPU Controller HTTP 测试

### 运行测试

```bash
# 运行所有 DAO 测试
go test -v ./internal/dao/...

# 运行所有 Service 测试
go test -v ./internal/service/...

# 运行所有 Controller 测试
go test -v ./internal/controller/v1/...

# 运行单个测试
go test -v ./internal/controller/v1/... -run TestHostController_Create
go test -v ./internal/controller/v1/... -run TestGPUController_Create
```

### 修复的问题
- GPU Controller 测试中 ID 转换错误：`string(rune(int(id)))` 改为 `fmt.Sprintf("/gpus/%d", int(id))`

---

## 2026-01-30: GPU 模块功能补充

### 问题
GPU 模块缺少 List、Update、Allocate、Release 功能。

### 解决方案

**新增 DAO 方法**:
- `List(page, pageSize)` - 分页获取GPU列表
- `GetByStatus(status)` - 按状态查询
- `Allocate(id, allocatedTo)` - 分配GPU
- `Release(id)` - 释放GPU

**新增 Service 方法**:
- `List` / `GetByStatus` / `Allocate` / `Release`
- Allocate/Release 包含状态校验逻辑

**新增 Controller 接口**:
- `List` / `Update` / `Allocate` / `Release`

### 新增 API 路由

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /admin/gpus | GPU列表 |
| PUT | /admin/gpus/:id | 更新GPU |
| POST | /admin/gpus/:id/allocate | 分配GPU |
| POST | /admin/gpus/:id/release | 释放GPU |

### 运行测试

```bash
go test -v ./internal/dao/... -run TestGPUDao
```

---
