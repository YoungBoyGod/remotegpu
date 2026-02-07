# SSH vs Web服务配置差异说明

## 核心差异

### 协议层级不同

**Web服务(HTTP/HTTPS):**
- 应用层协议(七层)
- nginx可以通过`Host`头识别域名
- 可以用一个端口(443)服务多个域名
- 支持基于域名的路由

**SSH:**
- 传输层协议(四层)
- 没有Host头的概念
- nginx的stream模块**不支持**根据域名路由
- 一个端口只能对应一个后端

---

## SSH访问方案对比

### 方案A: 直接访问企业公网IP(推荐)

**访问方式:**
```bash
ssh -p 2201 user@企业公网IP    # GPU1
ssh -p 2202 user@企业公网IP    # GPU2
ssh -p 2203 user@企业公网IP    # GPU3
```

**优点:**
- ✅ 最简单,不需要云服务器参与
- ✅ 性能最好,少一次转发
- ✅ 云服务器不需要开放2201-2400端口
- ✅ 延迟最低

**缺点:**
- ❌ 用户需要记住企业公网IP
- ❌ 需要记住端口号

**配置要求:**
- 只需要企业防火墙的端口转发规则
- 云服务器不需要配置

**网络流量路径:**
```
用户 → 企业公网IP:2201 → 防火墙转发 → GPU1内网IP:22
```

---

### 方案B: 云服务器端口中转

**访问方式:**
```bash
ssh -p 2201 user@gpu.domain.com    # GPU1
ssh -p 2202 user@gpu.domain.com    # GPU2
ssh -p 2203 user@gpu.domain.com    # GPU3
```

**优点:**
- ✅ 可以使用域名
- ✅ 统一入口
- ✅ 企业公网IP变更时只需修改云服务器配置

**缺点:**
- ❌ 需要云服务器开放2201-2400端口
- ❌ 多一次转发,性能略差
- ❌ 云服务器带宽消耗增加

**配置要求(云服务器):**
```nginx
stream {
    # GPU1 SSH
    server {
        listen 2201;
        proxy_pass 企业公网IP:2201;
    }

    # GPU2 SSH
    server {
        listen 2202;
        proxy_pass 企业公网IP:2202;
    }

    # GPU3-200 同理...
}
```

**网络流量路径:**
```
用户 → gpu.domain.com:2201 → 云服务器:2201 → 企业公网IP:2201 → 防火墙转发 → GPU1内网IP:22
```

---

### 方案C: SSH配置文件简化(配合方案A或B)

无论选择方案A还是B,都可以通过SSH配置文件简化访问。

**用户本地 `~/.ssh/config`:**
```
# GPU1
Host gpu1
    HostName 企业公网IP  # 或 gpu.domain.com
    Port 2201
    User your_username
    IdentityFile ~/.ssh/id_rsa

# GPU2
Host gpu2
    HostName 企业公网IP
    Port 2202
    User your_username
    IdentityFile ~/.ssh/id_rsa

# GPU3-200 同理...
```

**访问方式:**
```bash
ssh gpu1    # 自动使用配置的IP和端口
ssh gpu2
ssh gpu3
```

**批量生成配置脚本:**
```bash
#!/bin/bash
# generate_ssh_config.sh

ENTERPRISE_IP="企业公网IP"  # 或 gpu.domain.com
USERNAME="your_username"

for i in {1..200}; do
    SSH_PORT=$((2200 + i))
    cat >> ~/.ssh/config <<EOF
Host gpu${i}
    HostName ${ENTERPRISE_IP}
    Port ${SSH_PORT}
    User ${USERNAME}
    IdentityFile ~/.ssh/id_rsa

EOF
done

echo "SSH配置已生成到 ~/.ssh/config"
```

---

## Web服务配置(支持子域名)

Web服务可以使用子域名,因为HTTP协议支持Host头:

**云服务器nginx配置:**
```nginx
# GPU1 Jupyter
server {
    listen 443 ssl http2;
    server_name gpu1-jupyter.gpu.domain.com;

    ssl_certificate /etc/letsencrypt/live/gpu.domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/gpu.domain.com/privkey.pem;

    location / {
        proxy_pass http://企业公网IP:8001;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
    }
}

# GPU2 Jupyter
server {
    listen 443 ssl http2;
    server_name gpu2-jupyter.gpu.domain.com;

    ssl_certificate /etc/letsencrypt/live/gpu.domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/gpu.domain.com/privkey.pem;

    location / {
        proxy_pass http://企业公网IP:8002;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
    }
}
```

**访问方式:**
```bash
https://gpu1-jupyter.gpu.domain.com
https://gpu2-jupyter.gpu.domain.com
https://gpu1-tensorboard.gpu.domain.com
```

**网络流量路径:**
```
用户 → gpu1-jupyter.gpu.domain.com:443 → 云服务器nginx(SSL终结) → 企业公网IP:8001 → 防火墙转发 → GPU1内网IP:8888
```

---

## 总结对比表

| 特性 | SSH | Web服务 |
|------|-----|---------|
| **协议层级** | 四层(传输层) | 七层(应用层) |
| **nginx模块** | stream | http |
| **是否支持子域名** | ❌ 不支持 | ✅ 支持 |
| **一个端口服务多个后端** | ❌ 不可以 | ✅ 可以 |
| **Host头识别** | ❌ 无 | ✅ 有 |
| **推荐方案** | 直接访问企业IP | 通过云服务器代理 |

---

## 推荐配置方案

### SSH访问
- **推荐**: 方案A(直接访问企业公网IP)
- **原因**: 性能最优,配置最简单
- **优化**: 配合SSH配置文件简化访问

### Web服务访问
- **推荐**: 通过云服务器nginx代理
- **原因**:
  - 提供SSL加密
  - 支持子域名(用户体验好)
  - 统一证书管理
  - 企业IP变更时只需修改nginx配置

---

## 常见误区

### ❌ 错误配置示例

```nginx
# 这个配置是错误的!
stream {
    server {
        listen 22;
        server_name gpu1.gpu.domain.com;  # stream模块不支持server_name!
        proxy_pass 企业公网IP:2201;
    }
}
```

**错误原因:**
1. stream模块不支持`server_name`指令
2. stream是四层代理,无法识别域名
3. 监听22端口会和云服务器自己的SSH冲突

### ✅ 正确理解

- SSH无法像HTTP那样通过域名区分后端
- 每个SSH连接需要独立的端口
- 如果要使用域名,只能是统一域名+不同端口
