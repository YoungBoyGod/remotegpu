# 机器远程访问配置设计

## 1. 背景与目标
管理员需要为机器配置对外访问能力（开放端口、绑定域名），并可通过平台统一管理。该功能后续需对接 Nginx 反向代理，实现统一入口与访问控制。

本阶段目标：
- 管理端可配置对外访问域名与开放端口
- 支持保存与展示配置
- 预留后端 API 与 Nginx 对接方案

## 2. 功能范围
- 仅管理员侧可配置
- 可设置对外访问域名、端口、协议
- 支持目标端口（反向代理映射）
- 支持额外开放端口列表（如 VNC/TensorBoard）

## 3. 配置字段建议

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| enabled | bool | 是 | 是否启用对外访问 |
| protocol | string | 是 | tcp/http/https/ssh |
| public_domain | string | 是 | 对外访问域名 |
| public_port | int | 是 | 对外开放端口 |
| target_port | int | 否 | 目标服务端口（反向代理用） |
| extra_ports | string | 否 | 额外开放端口列表，逗号分隔 |
| remark | string | 否 | 备注说明 |

## 4. 自动生成逻辑（建议）

### 4.1 域名生成
- 默认域名格式：`<machine_slug>.<base_domain>`
- `machine_slug` 取机器 `hostname/name/id` 之一并做 slugify（小写、替换空格）
- `base_domain` 通过配置提供（如环境变量）

### 4.2 端口映射
- 默认根据 `target_port` 映射到 `public_port`
- 推荐的默认映射表（可配置）：

| target_port | public_port |
|-------------|-------------|
| 22 | 2222 |
| 80 | 8080 |
| 443 | 8443 |
| 3389 | 13389 |
| 5900 | 15900 |
| 6006 | 16006 |
| 8888 | 18888 |

- 若无映射规则，默认 `public_port = target_port`

### 4.3 协议默认端口
- http → 80
- https → 443
- ssh → 22

## 5. 前端表现（已实现）
- 机器详情页新增“远程访问配置”卡片
- 启用后自动生成域名与端口
- 可手动覆盖并保存
- 当前仅本地保存配置，后续对接后端

## 6. 数据模型建议

## 7. API 设计建议（Admin）

```sql
CREATE TABLE machine_remote_access (
  id SERIAL PRIMARY KEY,
  host_id VARCHAR(64) NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT FALSE,
  protocol VARCHAR(16) NOT NULL DEFAULT 'tcp',
  public_domain VARCHAR(255) NOT NULL,
  public_port INT NOT NULL,
  target_port INT,
  extra_ports VARCHAR(255),
  remark TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_remote_access_domain_port
  ON machine_remote_access(public_domain, public_port);
```

## 8. Nginx 对接方案

### 6.1 获取配置
```
GET /api/v1/admin/machines/:id/remote-access

Response:
{
  "enabled": true,
  "protocol": "https",
  "public_domain": "gpu-001.example.com",
  "public_port": 443,
  "target_port": 8888,
  "extra_ports": "6006,5901",
  "remark": "Jupyter + TensorBoard"
}
```

### 6.2 更新配置
```
PUT /api/v1/admin/machines/:id/remote-access

Request:
{
  "enabled": true,
  "protocol": "https",
  "public_domain": "gpu-001.example.com",
  "public_port": 443,
  "target_port": 8888,
  "extra_ports": "6006,5901",
  "remark": "Jupyter + TensorBoard"
}
```

### 6.3 触发 Nginx 刷新（可选）
```
POST /api/v1/admin/machines/:id/remote-access/sync
```

## 9. 安全建议

### 7.1 推荐方式
- 后端生成 Nginx 配置片段（server/location/upstream）
- 写入配置目录并触发 reload
- 或使用 Nginx Plus/OpenResty API 动态更新

### 7.2 示例配置

```nginx
server {
  listen 443 ssl;
  server_name gpu-001.example.com;

  location / {
    proxy_pass http://10.0.0.12:8888;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  }
}
```

### 7.3 校验与约束
- 域名唯一 + 端口唯一
- 端口可用性检查
- DNS 解析检查（可选）
- 白名单/访问控制（CIDR）

## 10. 下一步计划
- 管理员权限校验
- 启用 HTTPS 与证书管理
- 记录审计日志
- 对外访问默认关闭，需要显式启用

## 11. 前端默认配置项

建议在前端使用以下配置项：

```
VITE_REMOTE_ACCESS_BASE_DOMAIN=remote.example.com
```
- 增加后端 API 与持久化
- 对接 Nginx 反向代理与证书管理
- 在前端增加“同步配置”与状态展示
