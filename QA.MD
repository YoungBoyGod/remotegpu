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
