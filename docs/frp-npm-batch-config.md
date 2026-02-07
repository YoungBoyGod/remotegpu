# Nginx Proxy Manager 批量配置方案

## 概述

本文档介绍如何批量配置 NPM 代理主机，适用于需要为大量 GPU 机器（如 200 台）配置反向代理的场景。

### 批量配置的挑战

手动在 NPM 界面逐个添加代理主机非常耗时：
- 200 台机器 × 2 个服务（Jupyter + TensorBoard）= 400 个代理配置
- 每个配置需要 2-3 分钟 = 总计 13-20 小时

因此需要使用自动化方案。

---

## 方案对比

| 方案 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| **方案 1：NPM API** | 官方支持，安全可靠 | 需要处理认证，API 文档较少 | ⭐⭐⭐⭐ |
| **方案 2：数据库操作** | 速度快，批量插入 | 需要直接操作数据库，有风险 | ⭐⭐⭐ |
| **方案 3：配置文件导入** | 简单直接 | NPM 不支持配置导入功能 | ❌ 不可用 |

---

## 方案 1：使用 NPM API

### 1.1 API 认证

NPM 提供 REST API，需要先获取认证 token。

#### 获取 Token

```bash
# 登录获取 token
curl -X POST http://your-server-ip:81/api/tokens \
  -H "Content-Type: application/json" \
  -d '{
    "identity": "admin@example.com",
    "secret": "your-password"
  }'
```

**响应示例**：
```json
{
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires": "2026-02-08T12:00:00.000Z"
}
```

保存 token 到环境变量：
```bash
export NPM_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 1.2 添加代理主机

#### API 端点

```
POST http://your-server-ip:81/api/nginx/proxy-hosts
```

#### 请求示例

```bash
curl -X POST http://your-server-ip:81/api/nginx/proxy-hosts \
  -H "Authorization: Bearer $NPM_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "domain_names": ["gpu1-jupyter.gpu.domain.com"],
    "forward_scheme": "http",
    "forward_host": "127.0.0.1",
    "forward_port": 7001,
    "certificate_id": 0,
    "ssl_forced": false,
    "block_exploits": true,
    "allow_websocket_upgrade": true,
    "access_list_id": 0,
    "advanced_config": "",
    "enabled": true,
    "meta": {}
  }'
```

### 1.3 申请 SSL 证书

添加代理主机后，需要为其申请 SSL 证书。

#### 方法 A：使用通配符证书（推荐）

先申请一个通配符证书，然后所有代理主机共用。

```bash
# 1. 申请通配符证书（需要 DNS API 支持）
curl -X POST http://your-server-ip:81/api/nginx/certificates \
  -H "Authorization: Bearer $NPM_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "letsencrypt",
    "domain_names": ["*.gpu.domain.com"],
    "meta": {
      "letsencrypt_email": "admin@example.com",
      "dns_challenge": true,
      "dns_provider": "cloudflare",
      "dns_provider_credentials": "dns_cloudflare_api_token=your-token"
    }
  }'

# 2. 获取证书 ID（从响应中）
# 假设返回的证书 ID 为 1

# 3. 更新代理主机，使用证书
curl -X PUT http://your-server-ip:81/api/nginx/proxy-hosts/1 \
  -H "Authorization: Bearer $NPM_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "certificate_id": 1,
    "ssl_forced": true,
    "http2_support": true,
    "hsts_enabled": true,
    "hsts_subdomains": false
  }'
```

#### 方法 B：为每个域名申请证书

```bash
curl -X POST http://your-server-ip:81/api/nginx/certificates \
  -H "Authorization: Bearer $NPM_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "letsencrypt",
    "domain_names": ["gpu1-jupyter.gpu.domain.com"],
    "meta": {
      "letsencrypt_email": "admin@example.com",
      "dns_challenge": false
    }
  }'
```

### 1.4 批量配置脚本

创建 Python 脚本 `npm_batch_add.py`：

```python
#!/usr/bin/env python3
import requests
import json
import time

# 配置
NPM_URL = "http://your-server-ip:81"
NPM_EMAIL = "admin@example.com"
NPM_PASSWORD = "your-password"
DOMAIN_SUFFIX = "gpu.domain.com"
WILDCARD_CERT_ID = 1  # 通配符证书 ID

# 机器配置
MACHINES = [
    {"id": 1, "jupyter_port": 7001, "tensorboard_port": 7002},
    {"id": 2, "jupyter_port": 7003, "tensorboard_port": 7004},
    # ... 添加更多机器
]

def get_token():
    """获取认证 token"""
    response = requests.post(
        f"{NPM_URL}/api/tokens",
        json={"identity": NPM_EMAIL, "secret": NPM_PASSWORD}
    )
    response.raise_for_status()
    return response.json()["token"]

def add_proxy_host(token, domain, port, cert_id=0):
    """添加代理主机"""
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }

    data = {
        "domain_names": [domain],
        "forward_scheme": "http",
        "forward_host": "127.0.0.1",
        "forward_port": port,
        "certificate_id": cert_id,
        "ssl_forced": cert_id > 0,
        "http2_support": cert_id > 0,
        "hsts_enabled": cert_id > 0,
        "block_exploits": True,
        "allow_websocket_upgrade": True,
        "enabled": True,
        "meta": {}
    }

    response = requests.post(
        f"{NPM_URL}/api/nginx/proxy-hosts",
        headers=headers,
        json=data
    )
    response.raise_for_status()
    return response.json()

def main():
    print("正在获取认证 token...")
    token = get_token()
    print(f"Token 获取成功")

    total = len(MACHINES) * 2
    current = 0

    for machine in MACHINES:
        machine_id = machine["id"]

        # 添加 Jupyter 代理
        jupyter_domain = f"gpu{machine_id}-jupyter.{DOMAIN_SUFFIX}"
        print(f"[{current+1}/{total}] 添加 {jupyter_domain}...")
        try:
            add_proxy_host(token, jupyter_domain, machine["jupyter_port"], WILDCARD_CERT_ID)
            print(f"  ✓ 成功")
        except Exception as e:
            print(f"  ✗ 失败: {e}")
        current += 1
        time.sleep(0.5)  # 避免请求过快

        # 添加 TensorBoard 代理
        tb_domain = f"gpu{machine_id}-tensorboard.{DOMAIN_SUFFIX}"
        print(f"[{current+1}/{total}] 添加 {tb_domain}...")
        try:
            add_proxy_host(token, tb_domain, machine["tensorboard_port"], WILDCARD_CERT_ID)
            print(f"  ✓ 成功")
        except Exception as e:
            print(f"  ✗ 失败: {e}")
        current += 1
        time.sleep(0.5)

    print(f"\n完成！成功添加 {current} 个代理配置")

if __name__ == "__main__":
    main()
```

### 1.5 运行脚本

```bash
# 安装依赖
pip3 install requests

# 运行脚本
python3 npm_batch_add.py
```

---

## 方案 2：直接操作数据库

### 2.1 数据库位置

NPM 使用 SQLite 数据库存储配置：
```
/data/database.sqlite
```

在 Docker 容器中的路径：
```bash
docker exec -it npm ls -la /data/database.sqlite
```

### 2.2 备份数据库

**重要**：操作数据库前务必备份！

```bash
# 备份数据库
docker exec npm cp /data/database.sqlite /data/database.sqlite.backup

# 或从容器外备份
docker cp npm:/data/database.sqlite ./database.sqlite.backup
```

### 2.3 数据库结构

查看 `proxy_host` 表结构：

```sql
CREATE TABLE proxy_host (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_on DATETIME NOT NULL,
  modified_on DATETIME NOT NULL,
  owner_user_id INTEGER NOT NULL,
  domain_names TEXT NOT NULL,
  forward_scheme TEXT NOT NULL,
  forward_host TEXT NOT NULL,
  forward_port INTEGER NOT NULL,
  certificate_id INTEGER NOT NULL DEFAULT 0,
  ssl_forced INTEGER NOT NULL DEFAULT 0,
  hsts_enabled INTEGER NOT NULL DEFAULT 0,
  hsts_subdomains INTEGER NOT NULL DEFAULT 0,
  http2_support INTEGER NOT NULL DEFAULT 0,
  block_exploits INTEGER NOT NULL DEFAULT 0,
  caching_enabled INTEGER NOT NULL DEFAULT 0,
  allow_websocket_upgrade INTEGER NOT NULL DEFAULT 0,
  access_list_id INTEGER NOT NULL DEFAULT 0,
  advanced_config TEXT,
  enabled INTEGER NOT NULL DEFAULT 1,
  meta TEXT
);
```

### 2.4 批量插入脚本

创建 Python 脚本 `npm_db_batch_add.py`：

```python
#!/usr/bin/env python3
import sqlite3
import json
from datetime import datetime

# 配置
DB_PATH = "./database.sqlite"  # 从容器复制出来的数据库
DOMAIN_SUFFIX = "gpu.domain.com"
OWNER_USER_ID = 1  # 管理员用户 ID
CERT_ID = 1  # 通配符证书 ID

# 机器配置
MACHINES = [
    {"id": 1, "jupyter_port": 7001, "tensorboard_port": 7002},
    {"id": 2, "jupyter_port": 7003, "tensorboard_port": 7004},
    # ... 添加更多机器
]

def add_proxy_hosts():
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()

    now = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    for machine in MACHINES:
        machine_id = machine["id"]

        # Jupyter
        jupyter_domain = f'["gpu{machine_id}-jupyter.{DOMAIN_SUFFIX}"]'
        cursor.execute("""
            INSERT INTO proxy_host (
                created_on, modified_on, owner_user_id,
                domain_names, forward_scheme, forward_host, forward_port,
                certificate_id, ssl_forced, hsts_enabled, http2_support,
                block_exploits, allow_websocket_upgrade, enabled, meta
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """, (
            now, now, OWNER_USER_ID,
            jupyter_domain, "http", "127.0.0.1", machine["jupyter_port"],
            CERT_ID, 1, 1, 1,
            1, 1, 1, "{}"
        ))
        print(f"添加 gpu{machine_id}-jupyter")

        # TensorBoard
        tb_domain = f'["gpu{machine_id}-tensorboard.{DOMAIN_SUFFIX}"]'
        cursor.execute("""
            INSERT INTO proxy_host (
                created_on, modified_on, owner_user_id,
                domain_names, forward_scheme, forward_host, forward_port,
                certificate_id, ssl_forced, hsts_enabled, http2_support,
                block_exploits, allow_websocket_upgrade, enabled, meta
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """, (
            now, now, OWNER_USER_ID,
            tb_domain, "http", "127.0.0.1", machine["tensorboard_port"],
            CERT_ID, 1, 1, 1,
            1, 1, 1, "{}"
        ))
        print(f"添加 gpu{machine_id}-tensorboard")

    conn.commit()
    conn.close()
    print(f"\n完成！共添加 {len(MACHINES) * 2} 个代理配置")

if __name__ == "__main__":
    add_proxy_hosts()
```

### 2.5 执行步骤

```bash
# 1. 停止 NPM 容器
docker-compose stop npm

# 2. 复制数据库到本地
docker cp npm:/data/database.sqlite ./database.sqlite

# 3. 运行批量插入脚本
python3 npm_db_batch_add.py

# 4. 复制数据库回容器
docker cp ./database.sqlite npm:/data/database.sqlite

# 5. 启动 NPM 容器
docker-compose start npm

# 6. 检查日志
docker logs npm
```

---

## 生成机器配置列表

### 使用 Python 生成

```python
# 生成 200 台机器的配置
machines = []
base_port = 7001

for i in range(1, 201):
    machines.append({
        "id": i,
        "jupyter_port": base_port + (i-1) * 2,
        "tensorboard_port": base_port + (i-1) * 2 + 1
    })

# 输出为 JSON
import json
print(json.dumps(machines, indent=2))
```

### 使用 Bash 生成

```bash
# 生成机器配置（JSON 格式）
for i in {1..200}; do
  jupyter_port=$((7001 + (i-1) * 2))
  tb_port=$((7001 + (i-1) * 2 + 1))
  echo "{\"id\": $i, \"jupyter_port\": $jupyter_port, \"tensorboard_port\": $tb_port},"
done
```

---

## 验证批量配置

### 1. 检查代理主机数量

登录 NPM 管理界面，查看 **Proxy Hosts** 页面，应该看到 400 个代理配置。

### 2. 批量测试访问

创建测试脚本 `test_batch_access.sh`：

```bash
#!/bin/bash

DOMAIN_SUFFIX="gpu.domain.com"

for i in {1..200}; do
  jupyter_url="https://gpu${i}-jupyter.${DOMAIN_SUFFIX}"
  tb_url="https://gpu${i}-tensorboard.${DOMAIN_SUFFIX}"

  echo -n "测试 GPU${i} Jupyter... "
  if curl -s -o /dev/null -w "%{http_code}" "$jupyter_url" | grep -q "200\|302\|401"; then
    echo "✓"
  else
    echo "✗"
  fi

  echo -n "测试 GPU${i} TensorBoard... "
  if curl -s -o /dev/null -w "%{http_code}" "$tb_url" | grep -q "200\|302\|401"; then
    echo "✓"
  else
    echo "✗"
  fi
done
```

运行测试：
```bash
chmod +x test_batch_access.sh
./test_batch_access.sh
```

---

## 常见问题

### API 请求失败

**问题**：批量添加时部分请求失败

**解决方案**：
1. 增加请求间隔：`time.sleep(1)`
2. 添加重试机制
3. 检查 NPM 日志：`docker logs npm`

### 数据库操作后配置未生效

**问题**：数据库插入成功，但 NPM 界面看不到

**解决方案**：
```bash
# 重启 NPM 容器
docker-compose restart npm

# 检查数据库文件权限
docker exec npm ls -la /data/database.sqlite
```

### SSL 证书申请速率限制

**问题**：Let's Encrypt 限制每周申请次数

**解决方案**：
1. 使用通配符证书（推荐）
2. 分批申请，每批间隔 1 小时
3. 使用测试环境先验证：`letsencrypt_staging: true`

---

## 性能优化

### 并发请求

修改 API 脚本，使用多线程：

```python
from concurrent.futures import ThreadPoolExecutor

def add_all_proxies(token):
    with ThreadPoolExecutor(max_workers=10) as executor:
        futures = []
        for machine in MACHINES:
            futures.append(executor.submit(add_proxy_host, token, ...))

        for future in futures:
            try:
                future.result()
            except Exception as e:
                print(f"错误: {e}")
```

### 数据库批量插入

使用 `executemany` 提高插入速度：

```python
cursor.executemany("""
    INSERT INTO proxy_host (...) VALUES (?, ?, ...)
""", data_list)
```

---

## 下一步

批量配置完成后，继续阅读：

- **[NPM 完整实施指南](frp-npm-guide.md)** - 完整的 frp + NPM 方案
- **[frp 实施指南](frp-implementation-guide.md)** - frp 服务端和客户端配置

---

## 参考资源

- [NPM API 文档](https://nginxproxymanager.com/api/)
- [SQLite 官方文档](https://www.sqlite.org/docs.html)
- [Python Requests 库](https://requests.readthedocs.io/)
