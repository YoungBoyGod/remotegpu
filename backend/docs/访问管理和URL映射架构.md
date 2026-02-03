# 访问管理和URL映射架构

## 概述

本文档描述了RemoteGPU平台的访问管理架构,包括SSH登录、容器应用访问的URL映射规则和实现原理。

## 1. SSH登录信息展示格式

### 1.1 展示规范

每台机器的SSH登录信息按照以下格式展示:

```
SSH 登录信息
连接主机: 6hiflzwte2xnkroq.ssh.x-gpu.com
端口: 56376
用户: root    密码: kD2oTIydfOKmQadjma0gH5JyDApG0uWn
连接命令: ssh -p 56376 root@6hiflzwte2xnkroq.ssh.x-gpu.com
```

### 1.2 字段说明

- **连接主机**: 使用唯一的SSH域名,格式为 `{环境标识符}.ssh.x-gpu.com`
- **端口**: 高端口号(50000-60000范围),避免与标准端口冲突
- **用户**: SSH登录用户名(通常为root)
- **密码**: 32位随机生成的强密码
- **连接命令**: 完整的SSH连接命令,用户可直接复制使用

## 2. URL格式规则

### 2.1 三种访问方式

平台支持三种访问方式,每种使用不同的URL格式:

#### 2.1.1 SSH访问

**格式**: `{环境标识符}.ssh.x-gpu.com`

**示例**: `6hiflzwte2xnkroq.ssh.x-gpu.com`

**说明**:
- 用于SSH终端访问
- 需要配合端口号使用
- 通过SSH网关进行路由

#### 2.1.2 Jupyter访问

**格式**: `https://{环境标识符}-8888.container.x-gpu.com/lab/...`

**示例**: `https://6hiflzwte2xnkroq-8888.container.x-gpu.com/lab/workspaces/auto-H/tree/root`

**说明**:
- Jupyter默认使用8888端口
- 端口号编码在域名中
- 支持JupyterLab完整功能

#### 2.1.3 应用映射访问

**格式**: `https://{环境标识符}-{端口}.container.x-gpu.com/`

**示例**:
- Gradio应用: `https://iug28wj6h5n8a0me-7860.container.x-gpu.com/`
- Streamlit应用: `https://iug28wj6h5n8a0me-8501.container.x-gpu.com/`
- 自定义应用: `https://iug28wj6h5n8a0me-{自定义端口}.container.x-gpu.com/`

**说明**:
- 支持任意容器内应用的端口映射
- 端口号灵活配置
- 自动HTTPS加密

### 2.2 环境标识符生成规则

**格式**: 16位随机字符串(小写字母+数字)

**示例**: `6hiflzwte2xnkroq`, `iug28wj6h5n8a0me`

**生成方法**:
```go
// 使用UUID或随机字符串生成
identifier := generateRandomString(16) // 例如: 6hiflzwte2xnkroq
```

## 3. 实现架构

### 3.1 整体架构图

```
用户浏览器/SSH客户端
        ↓
DNS解析 (*.ssh.x-gpu.com / *.container.x-gpu.com)
        ↓
入口网关 (Nginx/Traefik)
        ↓
    路由解析
   ↙         ↘
SSH网关      容器访问网关
   ↓              ↓
目标容器SSH    目标容器应用端口
```

### 3.2 DNS配置

**泛域名解析**:
```
*.ssh.x-gpu.com        → SSH网关IP (例如: 192.168.10.100)
*.container.x-gpu.com  → 容器网关IP (例如: 192.168.10.101)
```

**SSL证书**:
- 使用泛域名SSL证书
- 支持Let's Encrypt自动续期
- 证书覆盖所有子域名

### 3.3 反向代理配置

#### 3.3.1 Nginx配置示例(容器访问)

```nginx
server {
    listen 443 ssl http2;
    server_name ~^(?<env_id>[a-z0-9]+)-(?<port>\d+)\.container\.x-gpu\.com$;

    ssl_certificate /etc/nginx/ssl/wildcard.container.x-gpu.com.crt;
    ssl_certificate_key /etc/nginx/ssl/wildcard.container.x-gpu.com.key;

    location / {
        # 根据环境ID查询容器IP
        set $container_ip '';
        set_by_lua_block $container_ip {
            local env_id = ngx.var.env_id
            local port = ngx.var.port
            -- 从Redis/数据库查询容器IP
            return query_container_ip(env_id)
        }

        proxy_pass http://$container_ip:$port;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

#### 3.3.2 SSH网关配置

使用SSH代理服务(如sshpiper)根据域名前缀路由到目标容器:

```yaml
# sshpiper配置
routes:
  - pattern: "*.ssh.x-gpu.com"
    handler: dynamic_route
    lookup: redis  # 从Redis查询环境ID对应的容器IP和端口
```

## 4. 技术优势

### 4.1 用户体验优势

✅ **友好的URL**: 无需记忆复杂的IP地址和端口号
✅ **统一访问入口**: 所有服务使用相同的域名体系
✅ **自动HTTPS**: 所有连接自动加密,保障安全
✅ **即时可用**: 新环境创建后立即获得访问域名

### 4.2 运维管理优势

✅ **动态扩展**: 支持无限数量的环境,无需手动配置
✅ **集中管理**: 统一的网关入口,便于监控和日志收集
✅ **故障隔离**: 单个环境故障不影响其他环境
✅ **灵活路由**: 支持基于域名的智能路由和负载均衡

### 4.3 安全性优势

✅ **端口隐藏**: 真实端口不对外暴露,降低攻击面
✅ **访问控制**: 可在网关层实现统一的访问控制策略
✅ **审计日志**: 所有访问经过网关,便于审计和追踪
✅ **DDoS防护**: 网关层可部署DDoS防护措施

## 5. 数据流程

### 5.1 环境创建流程

```
1. 用户创建环境
   ↓
2. 系统生成唯一环境标识符 (16位随机字符串)
   ↓
3. 创建容器/虚拟机
   ↓
4. 生成访问信息:
   - SSH域名: {标识符}.ssh.x-gpu.com
   - 应用域名: {标识符}-{端口}.container.x-gpu.com
   ↓
5. 将映射关系存储到Redis/数据库:
   - Key: 环境标识符
   - Value: 容器IP、端口、认证信息
   ↓
6. 返回访问信息给用户
```

### 5.2 访问请求流程

```
1. 用户访问 https://6hiflzwte2xnkroq-8888.container.x-gpu.com
   ↓
2. DNS解析到容器网关IP
   ↓
3. Nginx解析域名:
   - 提取环境标识符: 6hiflzwte2xnkroq
   - 提取端口号: 8888
   ↓
4. 从Redis查询容器IP (例如: 192.168.100.50)
   ↓
5. 反向代理到 http://192.168.100.50:8888
   ↓
6. 返回响应给用户
```

## 6. 安全考虑

### 6.1 认证与授权

- **SSH密钥认证**: 推荐使用SSH密钥替代密码认证
- **JWT令牌**: 容器应用访问可使用JWT进行身份验证
- **访问控制列表**: 在网关层实现IP白名单/黑名单
- **多因素认证**: 对敏感操作启用MFA

### 6.2 网络隔离

- **容器网络隔离**: 不同用户的容器使用独立的网络命名空间
- **防火墙规则**: 限制容器只能访问必要的外部资源
- **内网访问**: 容器间通信仅限内网,不对外暴露

### 6.3 数据加密

- **传输加密**: 所有HTTP流量强制使用HTTPS
- **SSH加密**: SSH连接使用强加密算法
- **密码存储**: 使用bcrypt/argon2加密存储密码

### 6.4 审计与监控

- **访问日志**: 记录所有访问请求(时间、来源IP、目标环境)
- **异常检测**: 监控异常访问模式(频繁失败、异常流量)
- **告警机制**: 异常情况及时通知管理员

## 7. 实施建议

### 7.1 基础设施准备

1. **域名准备**
   - 注册主域名 `x-gpu.com`
   - 配置DNS泛域名解析
   - 申请泛域名SSL证书

2. **网关服务器**
   - SSH网关: 2核4GB内存(可根据并发量调整)
   - 容器网关: 4核8GB内存(推荐使用Nginx + OpenResty)
   - 高可用部署: 建议使用主备或负载均衡

3. **存储服务**
   - Redis: 用于存储环境ID到容器IP的映射关系
   - 数据库: 存储完整的环境配置和访问日志

### 7.2 部署步骤

1. **配置DNS**
   ```bash
   # 添加A记录
   *.ssh.x-gpu.com        A    192.168.10.100
   *.container.x-gpu.com  A    192.168.10.101
   ```

2. **部署Nginx网关**
   ```bash
   # 安装OpenResty(包含Lua支持)
   apt-get install openresty

   # 配置Nginx
   cp nginx.conf /etc/openresty/nginx.conf
   systemctl restart openresty
   ```

3. **配置Redis映射**
   ```bash
   # 环境创建时写入映射
   redis-cli SET "env:6hiflzwte2xnkroq:ip" "192.168.100.50"
   redis-cli SET "env:6hiflzwte2xnkroq:ssh_port" "56376"
   ```

### 7.3 监控与维护

- **性能监控**: 使用Prometheus监控网关性能
- **日志分析**: 使用ELK Stack分析访问日志
- **定期备份**: 备份Redis数据和环境配置
- **证书更新**: 自动化SSL证书续期

## 8. 总结

本架构通过DNS泛域名解析和反向代理技术,实现了:

✅ **用户友好**: 简洁易记的访问URL
✅ **自动化**: 环境创建即可访问,无需手动配置
✅ **安全可靠**: 统一的安全策略和访问控制
✅ **易于扩展**: 支持无限数量的环境和应用

该架构特别适合多租户GPU云平台场景,能够为每个用户环境提供独立、安全、易用的访问方式。

---

**文档版本**: v1.0
**最后更新**: 2026-02-02
**维护者**: RemoteGPU开发团队
