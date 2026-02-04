# CodeX SSHKey / 审计 / 镜像 Review（2026-02-04）

## Reviewed Files

- `internal/controller/v1/customer/sshkey_controller.go`
- `internal/service/sshkey/sshkey_service.go`
- `internal/controller/v1/ops/audit_controller.go`
- `internal/service/audit/audit_service.go`
- `internal/controller/v1/ops/image_controller.go`

## 发现的问题（原实现）

1. **SSHKey Controller 使用类型断言**
   - 原实现 `userID.(uint)` 可能在缺失或类型不匹配时 panic。
   - 影响：未授权或异常请求会直接导致服务崩溃。

2. **SSH 公钥解析不够稳健**
   - 原实现手动解析并限定 key 类型，容易遗漏合法类型或注释格式。
   - 影响：合法公钥可能被误判为无效，兼容性差。

3. **审计 / 镜像分页参数无保护**
   - 原实现未限制 `page`、`page_size`，可能出现负 offset 或极大分页。
   - 影响：数据库查询异常或性能风险。

4. **审计详情序列化忽略错误**
   - 原实现忽略 JSON 序列化错误，写入数据可能不完整。
   - 影响：审计字段丢失难以排查。

## 优化与改进（当前实现）

1. **安全获取 userID**
   - 使用 `GetUint` 和统一的 401 处理，避免 panic。

2. **使用 `ssh.ParseAuthorizedKey` 解析公钥**
   - 由标准库解析，自动支持更多合法格式。
   - 指纹使用 `FingerprintSHA256` 生成，更标准。

3. **分页参数归一化**
   - `page < 1`、`page_size` 超出范围时自动回退默认值。

4. **审计详情显式序列化**
   - `json.Marshal` 失败直接返回错误，避免悄悄写入失败。

## 为什么这样更好

- **更稳定**：避免 panic，异常请求不会拖垮服务。
- **更兼容**：标准库解析公钥，减少误判。
- **更可控**：分页有上限，数据库压力可预测。
- **更可追踪**：审计详情序列化可失败即报错。

## 备注

- 关键修改已在代码中添加 `CodeX 2026-02-04` 注释。
