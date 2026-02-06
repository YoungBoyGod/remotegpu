# SSH 密码加密迁移指南

## 概述

此迁移脚本用于将数据库中的明文 SSH 密码加密存储，修复 P0 安全问题。

## 前置条件

1. **设置加密密钥**（必须）
   ```bash
   export ENCRYPTION_KEY="your-32-byte-encryption-key!!"
   ```
   - 密钥必须是 32 字节（用于 AES-256）
   - 生产环境请使用强随机密钥
   - **重要**：密钥丢失将无法解密已加密的密码

2. **备份数据库**（强烈建议）
   ```bash
   pg_dump -h localhost -U postgres remotegpu > backup_$(date +%Y%m%d_%H%M%S).sql
   ```

## 使用方法

### 1. 试运行（推荐先执行）

检查需要迁移的记录，不会修改数据库：

```bash
cd backend
go run cmd/migrate_ssh_passwords/main.go --dry-run
```

输出示例：
```
=== SSH 密码加密迁移工具 ===
总记录数: 10
需要迁移: 8
模式: 试运行（不会修改数据）

[1/10] 迁移 node-01 ... OK (试运行)
[2/10] 迁移 node-02 ... OK (试运行)
...
```

### 2. 实际执行

确认无误后，执行实际迁移：

```bash
cd backend
go run cmd/migrate_ssh_passwords/main.go
```

系统会要求确认：
```
确认执行迁移？这将修改数据库中的 SSH 密码字段。(yes/no): yes
```

### 3. 强制执行（跳过确认）

```bash
cd backend
go run cmd/migrate_ssh_passwords/main.go --force
```

## 迁移逻辑

1. **读取所有 hosts 记录**
2. **跳过以下情况**：
   - SSH 密码为空
   - 已经加密的密码（通过长度和字符集判断）
3. **加密处理**：
   - 使用 AES-256-GCM 加密
   - Base64 编码存储
4. **更新数据库**

## 验证迁移

迁移完成后，验证加密是否成功：

```sql
-- 查看加密后的密码（应该是长字符串）
SELECT id, ssh_username,
       LENGTH(ssh_password) as password_length,
       LEFT(ssh_password, 20) as password_preview
FROM hosts
WHERE ssh_password IS NOT NULL AND ssh_password != '';
```

加密后的密码特征：
- 长度通常 > 40 字符
- 只包含 base64 字符（A-Z, a-z, 0-9, +, /, =）

## 测试连接

迁移后测试 SSH 连接是否正常：

1. 通过管理后台添加新机器
2. 触发机器信息采集
3. 检查日志确认 SSH 连接成功

## 故障排查

### 问题 1：ENCRYPTION_KEY 未设置

```
错误：必须设置 ENCRYPTION_KEY 环境变量（32 字节）
```

**解决**：设置环境变量
```bash
export ENCRYPTION_KEY="your-32-byte-encryption-key!!"
```

### 问题 2：密钥长度错误

```
错误：ENCRYPTION_KEY 必须是 32 字节
```

**解决**：确保密钥正好是 32 字节（32 个字符）

### 问题 3：迁移失败

```
[5/10] 迁移 node-05 ... 失败: failed to encrypt SSH password: ...
```

**解决**：
1. 检查数据库连接
2. 检查密码字段是否包含特殊字符
3. 查看详细错误日志

### 问题 4：连接失败

迁移后 SSH 连接失败：

**可能原因**：
1. 加密密钥不一致（运行时使用的密钥与迁移时不同）
2. 密码本身错误

**解决**：
1. 确认 ENCRYPTION_KEY 环境变量一致
2. 重新设置正确的 SSH 密码

## 回滚方案

如果迁移后出现问题，可以：

### 方案 1：从备份恢复

```bash
psql -h localhost -U postgres remotegpu < backup_20260206_120000.sql
```

### 方案 2：手动重置密码

通过管理后台重新设置受影响机器的 SSH 密码。

## 注意事项

1. **密钥管理**
   - 生产环境必须使用强随机密钥
   - 密钥应存储在安全的密钥管理系统（如 AWS KMS、HashiCorp Vault）
   - 定期轮换密钥（需要重新加密所有密码）

2. **性能考虑**
   - 大量记录（>1000）建议分批迁移
   - 可以修改脚本添加批处理逻辑

3. **并发安全**
   - 迁移期间避免创建新机器
   - 建议在维护窗口执行

4. **审计**
   - 记录迁移操作日志
   - 保留迁移前后的数据快照

## 后续步骤

迁移完成后：

1. ✅ 验证所有机器的 SSH 连接正常
2. ✅ 更新运维文档，说明密码已加密
3. ✅ 配置密钥管理流程
4. ✅ 删除备份文件（如果确认无问题）
5. ✅ 监控系统日志，确认无解密错误

## 技术细节

### 加密算法

- **算法**：AES-256-GCM
- **密钥长度**：256 位（32 字节）
- **模式**：GCM（Galois/Counter Mode）
- **认证**：内置消息认证
- **编码**：Base64

### 数据格式

加密后的数据格式：
```
[nonce (12 bytes)][encrypted data][auth tag (16 bytes)]
```

Base64 编码后存储到数据库。

## 相关文件

- 迁移脚本：`backend/cmd/migrate_ssh_passwords/main.go`
- 加密工具：`backend/pkg/crypto/aes.go`
- 机器服务：`backend/internal/service/machine/machine_service.go`
- 注册服务：`backend/internal/service/machine/enrollment_service.go`
