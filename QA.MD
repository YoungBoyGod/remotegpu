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
S3_ACCESS_KEY="your_key" \
S3_SECRET_KEY="your_secret" \
S3_BUCKET="remotegpu" \
go test -v ./pkg/storage/... -run TestS3Storage_Upload
```

### 相关文件
- `backend/pkg/storage/storage_test.go` - 测试文件
- `backend/pkg/storage/local.go` - 本地存储实现
- `backend/pkg/storage/s3.go` - S3存储实现

---
