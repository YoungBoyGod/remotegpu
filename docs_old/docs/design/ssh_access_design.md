# SSH 公网访问技术设计方案

> 开发环境 SSH 远程访问技术实现方案
>
> 创建日期：2026-01-26

---

## 目录

1. [整体架构](#1-整体架构)
2. [SSH 接入方案对比](#2-ssh-接入方案对比)
3. [推荐方案：容器化 SSH 服务](#3-推荐方案容器化-ssh-服务)
4. [用户认证方案](#4-用户认证方案)
5. [端口管理方案](#5-端口管理方案)
6. [网络架构设计](#6-网络架构设计)
7. [安全机制](#7-安全机制)
8. [实现流程](#8-实现流程)
9. [技术实现细节](#9-技术实现细节)

---

## 1. 整体架构

### 1.1 系统架构图

```
┌──────────────────────────────────────────────────────────────┐
│                        用户终端                                 │
│              ssh user123@ssh.example.com:30001                │
└────────────────────────┬─────────────────────────────────────┘
                         │
                         │ 公网访问
                         ▼
┌──────────────────────────────────────────────────────────────┐
│                    负载均衡器 / 入口节点                          │
│                  (Load Balancer / Gateway)                    │
│                     ssh.example.com                           │
└────────────────────────┬─────────────────────────────────────┘
                         │
                         │ 端口映射 (30001 -> Pod:22)
                         ▼
┌──────────────────────────────────────────────────────────────┐
│                    Kubernetes 集群                             │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  开发环境 Pod (dev-env-user123)                           │ │
│  │  ┌─────────────────────────────────────────────────────┐│ │
│  │  │  SSH Server (Port 22)                               ││ │
│  │  │  - 用户: user123                                     ││ │
│  │  │  - 密码: auto-generated                             ││ │
│  │  ├─────────────────────────────────────────────────────┤│ │
│  │  │  JupyterLab (Port 8888)                             ││ │
│  │  ├─────────────────────────────────────────────────────┤│ │
│  │  │  工作目录                                            ││ │
│  │  │  - /gemini/code/      (持久化存储)                  ││ │
│  │  │  - /gemini/data-*/    (数据集挂载)                  ││ │
│  │  │  - /gemini/output/    (输出目录)                    ││ │
│  │  └─────────────────────────────────────────────────────┘│ │
│  └─────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────────────┐
│                  持久化存储 (PVC)                               │
│              - 用户代码                                         │
│              - 数据集                                           │
│              - 模型                                             │
└──────────────────────────────────────────────────────────────┘
```

### 1.2 核心组件

1. **API 服务**：管理开发环境的创建/销毁/配置
2. **SSH 网关**：负责 SSH 连接的路由和转发
3. **Kubernetes 集群**：运行开发环境容器
4. **持久化存储**：存储用户数据
5. **数据库**：存储环境配置、端口映射、认证信息

---

## 2. SSH 接入方案对比

### 方案 A：每个容器独立 SSH 端口（推荐）

**架构：**
```
用户 -> LoadBalancer:30001 -> Pod:22 (dev-env-user123)
用户 -> LoadBalancer:30002 -> Pod:22 (dev-env-user456)
```

**优点：**
- ✅ 实现简单，每个环境完全隔离
- ✅ 用户体验好：`ssh user@host:port`
- ✅ 容器内标准 SSH 服务器（Port 22）
- ✅ 易于调试和维护

**缺点：**
- ❌ 需要大量公网端口（每个环境一个）
- ❌ 端口资源有限（NodePort 范围：30000-32767）

**适用场景：** 中小规模（< 1000 并发环境）

---

### 方案 B：SSH 网关 + 用户名路由

**架构：**
```
用户 -> ssh user123@gateway.com:22 -> SSH Gateway -> Pod (dev-env-user123)
用户 -> ssh user456@gateway.com:22 -> SSH Gateway -> Pod (dev-env-user456)
```

**实现方式：**
- 统一 SSH 网关（Port 22）
- 根据用户名路由到不同容器
- 使用 ProxyJump 或 SSH Proxy

**优点：**
- ✅ 节省端口资源（只需一个端口）
- ✅ 可扩展性好
- ✅ 统一入口，便于安全控制

**缺点：**
- ❌ 实现复杂，需要自定义 SSH 网关
- ❌ 需要处理 SSH 连接转发
- ❌ 网关成为单点（需要高可用）

**适用场景：** 大规模（> 1000 并发环境）

---

### 方案 C：共享主机 SSH + 容器隔离

**架构：**
```
用户 -> ssh user123@host1.com:22 -> 宿主机 SSH -> nsenter 进入容器
```

**实现方式：**
- 宿主机运行 SSH 服务器
- 为每个用户创建 Linux 用户
- 登录后通过 nsenter 进入对应容器

**优点：**
- ✅ 节省端口资源
- ✅ 标准 Linux 用户管理

**缺点：**
- ❌ 宿主机用户管理复杂
- ❌ 安全风险高（用户在宿主机上）
- ❌ 难以实现完全隔离
- ❌ 不适合 Kubernetes 环境

**适用场景：** 不推荐（安全风险）

---

### 方案 D：Websocket SSH（Web 终端）

**架构：**
```
用户浏览器 -> WebSocket -> API Gateway -> Kubernetes Pod
```

**实现方式：**
- 前端使用 xterm.js
- 后端 WebSocket 连接到容器

**优点：**
- ✅ 无需客户端配置
- ✅ 跨平台，浏览器即可使用
- ✅ 易于集成到 Web 界面

**缺点：**
- ❌ 不支持标准 SSH 客户端
- ❌ 不支持 IDE 远程开发（VSCode Remote-SSH）
- ❌ 功能受限（无法使用 scp、rsync 等）

**适用场景：** 作为补充方案，与方案 A 或 B 配合使用

---

## 3. 推荐方案：容器化 SSH 服务

### 3.1 方案选择

**推荐：方案 A（独立端口）+ 方案 D（Web 终端）**

- MVP 阶段：使用方案 A，快速实现
- 扩展阶段：迁移到方案 B（SSH 网关）
- 补充方案：同时提供方案 D（Web 终端）

### 3.2 实现架构

```
┌─────────────────────────────────────────────────────────────┐
│                      用户访问方式                              │
├─────────────────────────────────────────────────────────────┤
│  1. SSH 客户端: ssh user123@ssh.example.com:30001          │
│  2. Web 终端:   https://platform.com/terminal/env-123       │
│  3. VSCode:     Remote-SSH 配置                              │
│  4. SFTP:       sftp://ssh.example.com:30001                │
└─────────────────────────────────────────────────────────────┘
```

---

## 4. 用户认证方案

### 4.1 认证方式

#### 方式 1：密码认证（推荐作为默认）

**优点：**
- 简单易用，用户无需配置
- 适合临时环境

**实现：**
```bash
# 容器启动时自动生成密码
PASSWORD=$(openssl rand -base64 12)
echo "user123:${PASSWORD}" | chpasswd

# 存储到数据库
INSERT INTO ssh_credentials (env_id, username, password, created_at)
VALUES ('env-123', 'user123', '${PASSWORD}', NOW());
```

**安全机制：**
- 失败 5 次锁定 30 分钟
- 密码复杂度要求
- 定期强制修改密码

#### 方式 2：SSH 密钥认证（推荐作为可选）

**优点：**
- 更安全
- 适合长期使用的环境
- IDE 集成友好

**实现：**
```bash
# 用户上传公钥
POST /api/ssh-keys
{
  "name": "My Laptop",
  "public_key": "ssh-rsa AAAAB3NzaC1yc2E..."
}

# 容器启动时挂载公钥
mkdir -p /home/user123/.ssh
echo "${PUBLIC_KEY}" > /home/user123/.ssh/authorized_keys
chmod 600 /home/user123/.ssh/authorized_keys
```

### 4.2 用户管理

#### 选项 1：固定用户名（推荐）

```bash
# 每个容器使用固定用户名
USERNAME="developer"

# 优点：
# - 简单统一
# - Dockerfile 预配置
# - 不需要动态创建用户
```

#### 选项 2：动态用户名

```bash
# 基于客户 ID 生成用户名
USERNAME="user${CUSTOMER_ID}"  # 例如: user12345

# 优点：
# - 更好的隔离
# - 便于审计
# - 个性化

# 缺点：
# - 需要动态创建 Linux 用户
# - Dockerfile 需要支持
```

**推荐：选项 1（固定用户名 `developer`）**

---

## 5. 端口管理方案

### 5.1 端口分配策略

#### 方案 A：随机端口分配

```go
// 从端口池中随机分配一个可用端口
func AllocatePort() (int, error) {
    minPort := 30000
    maxPort := 32767

    for {
        port := rand.Intn(maxPort-minPort) + minPort
        if !isPortInUse(port) {
            return port, nil
        }
    }
}
```

#### 方案 B：顺序端口分配

```go
// 维护一个端口池，顺序分配
type PortPool struct {
    nextPort int
    usedPorts map[int]bool
}

func (p *PortPool) Allocate() int {
    for {
        port := p.nextPort
        p.nextPort++
        if p.nextPort > 32767 {
            p.nextPort = 30000
        }

        if !p.usedPorts[port] {
            p.usedPorts[port] = true
            return port
        }
    }
}
```

**推荐：方案 B（顺序分配）+ 端口回收机制**

### 5.2 端口映射配置

#### Kubernetes Service 配置

```yaml
apiVersion: v1
kind: Service
metadata:
  name: dev-env-user123-ssh
  namespace: dev-environments
spec:
  type: NodePort
  selector:
    app: dev-env
    env-id: env-123
  ports:
  - name: ssh
    protocol: TCP
    port: 22          # Service 内部端口
    targetPort: 22    # 容器端口
    nodePort: 30001   # 外部访问端口
```

### 5.3 端口映射存储

```sql
-- 端口映射表
CREATE TABLE port_mappings (
    id BIGSERIAL PRIMARY KEY,
    env_id VARCHAR(64) NOT NULL,
    customer_id BIGINT NOT NULL,
    external_port INT NOT NULL UNIQUE,
    internal_port INT NOT NULL DEFAULT 22,
    service_name VARCHAR(128),
    status VARCHAR(20) DEFAULT 'active',  -- active, released
    allocated_at TIMESTAMP DEFAULT NOW(),
    released_at TIMESTAMP,
    INDEX idx_env_id (env_id),
    INDEX idx_customer_id (customer_id),
    INDEX idx_external_port (external_port)
);
```

---

## 6. 网络架构设计

### 6.1 网络拓扑

```
Internet
    │
    ▼
[负载均衡器 / 边界路由器]
    │
    ├─ Port 30001 -> Node1:30001
    ├─ Port 30002 -> Node2:30002
    └─ Port 30003 -> Node3:30003
    │
    ▼
[Kubernetes 集群]
    │
    ├─ Node1
    │   └─ Pod (dev-env-user123) :22
    │
    ├─ Node2
    │   └─ Pod (dev-env-user456) :22
    │
    └─ Node3
        └─ Pod (dev-env-user789) :22
```

### 6.2 DNS 配置

```
# 统一 SSH 入口域名
ssh.example.com  ->  负载均衡器 IP

# 用户连接方式
ssh developer@ssh.example.com -p 30001
```

### 6.3 防火墙规则

```bash
# 开放 NodePort 端口范围
iptables -A INPUT -p tcp --dport 30000:32767 -j ACCEPT

# 限制来源 IP（可选）
iptables -A INPUT -p tcp --dport 30000:32767 -s <ALLOWED_IP> -j ACCEPT
iptables -A INPUT -p tcp --dport 30000:32767 -j DROP
```

---

## 7. 安全机制

### 7.1 认证安全

#### 密码策略

```go
type PasswordPolicy struct {
    MinLength      int  // 最小长度：12
    RequireUpper   bool // 需要大写字母
    RequireLower   bool // 需要小写字母
    RequireDigit   bool // 需要数字
    RequireSpecial bool // 需要特殊字符
}

// 生成强密码
func GeneratePassword() string {
    const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
    password := make([]byte, 12)
    for i := range password {
        password[i] = chars[rand.Intn(len(chars))]
    }
    return string(password)
}
```

#### 登录失败限制

```go
// 记录登录失败
type LoginAttempt struct {
    EnvID      string
    IPAddress  string
    Timestamp  time.Time
    Success    bool
}

// 检查是否被锁定
func IsLocked(envID string, ip string) bool {
    attempts := getFailedAttempts(envID, ip, 30*time.Minute)
    return len(attempts) >= 5
}

// SSH 配置
// /etc/ssh/sshd_config
MaxAuthTries 5
LoginGraceTime 2m
```

### 7.2 网络安全

#### IP 白名单（可选）

```go
// 配置允许访问的 IP 段
type IPWhitelist struct {
    CustomerID int64
    IPRanges   []string // ["192.168.1.0/24", "10.0.0.1/32"]
}

// 在 iptables 中配置
func ApplyIPWhitelist(envID string, whitelist []string) error {
    for _, ipRange := range whitelist {
        cmd := fmt.Sprintf(
            "iptables -A INPUT -p tcp --dport %d -s %s -j ACCEPT",
            port, ipRange,
        )
        exec.Command("sh", "-c", cmd).Run()
    }
    return nil
}
```

#### DDoS 防护

```yaml
# 使用 fail2ban
[sshd]
enabled = true
port = ssh
filter = sshd
logpath = /var/log/auth.log
maxretry = 5
bantime = 1800
```

### 7.3 容器安全

#### 资源限制

```yaml
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: dev-env
    resources:
      limits:
        cpu: "4"
        memory: "16Gi"
        nvidia.com/gpu: "1"
      requests:
        cpu: "1"
        memory: "4Gi"
```

#### 安全上下文

```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 1000
  fsGroup: 1000
  capabilities:
    drop:
    - ALL
    add:
    - NET_BIND_SERVICE  # 如果需要绑定特权端口
```

---

## 8. 实现流程

### 8.1 开发环境创建流程

```
┌────────────────────────────────────────────────────────────┐
│ 1. 用户请求创建开发环境                                       │
│    POST /api/environments                                  │
│    {                                                       │
│      "name": "My Dev Env",                                 │
│      "image": "ubuntu20.04-pytorch",                       │
│      "resources": {"cpu": 4, "memory": 16, "gpu": 1}       │
│    }                                                       │
└────────────────┬───────────────────────────────────────────┘
                 │
                 ▼
┌────────────────────────────────────────────────────────────┐
│ 2. 后端服务处理                                              │
│    a. 验证用户配额                                           │
│    b. 分配端口（30001）                                      │
│    c. 生成 SSH 密码                                         │
│    d. 创建数据库记录                                         │
└────────────────┬───────────────────────────────────────────┘
                 │
                 ▼
┌────────────────────────────────────────────────────────────┐
│ 3. 创建 Kubernetes 资源                                      │
│    a. 创建 PVC（持久化存储）                                 │
│    b. 创建 Pod（开发环境容器）                               │
│    c. 创建 Service（NodePort）                              │
│    d. 等待 Pod Ready                                        │
└────────────────┬───────────────────────────────────────────┘
                 │
                 ▼
┌────────────────────────────────────────────────────────────┐
│ 4. 容器初始化                                                │
│    a. 启动 SSH 服务                                         │
│    b. 启动 JupyterLab                                       │
│    c. 挂载数据集                                            │
│    d. 配置环境变量                                           │
└────────────────┬───────────────────────────────────────────┘
                 │
                 ▼
┌────────────────────────────────────────────────────────────┐
│ 5. 返回连接信息给用户                                         │
│    {                                                       │
│      "env_id": "env-123",                                  │
│      "ssh_host": "ssh.example.com",                        │
│      "ssh_port": 30001,                                    │
│      "ssh_username": "developer",                          │
│      "ssh_password": "xK9$mP2#qL5@",                       │
│      "jupyter_url": "https://platform.com/jupyter/env-123",│
│      "status": "running"                                   │
│    }                                                       │
└────────────────────────────────────────────────────────────┘
```

### 8.2 连接信息展示

```
┌─────────────────────────────────────────────────────────────┐
│               开发环境 SSH 连接信息                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  SSH 连接命令：                                              │
│  ┌─────────────────────────────────────────────────────────┐│
│  │ ssh developer@ssh.example.com -p 30001                 ││
│  └─────────────────────────────────────────────────────────┘│
│                                                             │
│  用户名：developer                                          │
│  密码：xK9$mP2#qL5@   [复制] [重置密码]                     │
│                                                             │
│  ───────────────────────────────────────────────────────    │
│                                                             │
│  VSCode Remote-SSH 配置：                                   │
│  ┌─────────────────────────────────────────────────────────┐│
│  │ Host myenv                                             ││
│  │   HostName ssh.example.com                             ││
│  │   Port 30001                                           ││
│  │   User developer                                       ││
│  └─────────────────────────────────────────────────────────┘│
│  [复制配置] [查看教程]                                       │
│                                                             │
│  ───────────────────────────────────────────────────────    │
│                                                             │
│  SFTP 配置：                                                │
│  地址：sftp://ssh.example.com:30001                         │
│  用户名：developer                                          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 9. 技术实现细节

### 9.1 Dockerfile 设计

```dockerfile
FROM ubuntu:20.04

# 安装 SSH 服务器
RUN apt-get update && \
    apt-get install -y openssh-server sudo && \
    mkdir /var/run/sshd

# 创建用户
RUN useradd -m -s /bin/bash developer && \
    echo "developer ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# 配置 SSH
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin no/' /etc/ssh/sshd_config && \
    sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config && \
    sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/' /etc/ssh/sshd_config

# 安装 Python 和 JupyterLab
RUN apt-get install -y python3 python3-pip && \
    pip3 install jupyterlab torch torchvision

# 创建工作目录
RUN mkdir -p /gemini/code /gemini/output && \
    chown -R developer:developer /gemini

# 启动脚本
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 22 8888

ENTRYPOINT ["/entrypoint.sh"]
```

### 9.2 容器启动脚本

```bash
#!/bin/bash
# entrypoint.sh

set -e

# 从环境变量读取密码
if [ -n "$SSH_PASSWORD" ]; then
    echo "developer:$SSH_PASSWORD" | chpasswd
    echo "SSH password set successfully"
fi

# 配置 SSH 公钥（如果提供）
if [ -n "$SSH_PUBLIC_KEY" ]; then
    mkdir -p /home/developer/.ssh
    echo "$SSH_PUBLIC_KEY" > /home/developer/.ssh/authorized_keys
    chmod 700 /home/developer/.ssh
    chmod 600 /home/developer/.ssh/authorized_keys
    chown -R developer:developer /home/developer/.ssh
    echo "SSH public key configured"
fi

# 启动 SSH 服务
/usr/sbin/sshd -D &
echo "SSH service started on port 22"

# 启动 JupyterLab（可选）
if [ "$ENABLE_JUPYTER" = "true" ]; then
    su - developer -c "jupyter lab --ip=0.0.0.0 --port=8888 --no-browser --notebook-dir=/gemini/code" &
    echo "JupyterLab started on port 8888"
fi

# 保持容器运行
wait
```

### 9.3 Kubernetes 资源创建

```go
package k8s

import (
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
)

// 创建开发环境 Pod
func CreateDevEnvironmentPod(envID, image, password string, resources ResourceConfig) error {
    pod := &corev1.Pod{
        ObjectMeta: metav1.ObjectMeta{
            Name:      fmt.Sprintf("dev-env-%s", envID),
            Namespace: "dev-environments",
            Labels: map[string]string{
                "app":    "dev-env",
                "env-id": envID,
            },
        },
        Spec: corev1.PodSpec{
            Containers: []corev1.Container{
                {
                    Name:  "dev-env",
                    Image: image,
                    Env: []corev1.EnvVar{
                        {Name: "SSH_PASSWORD", Value: password},
                        {Name: "ENABLE_JUPYTER", Value: "true"},
                    },
                    Ports: []corev1.ContainerPort{
                        {Name: "ssh", ContainerPort: 22},
                        {Name: "jupyter", ContainerPort: 8888},
                    },
                    Resources: corev1.ResourceRequirements{
                        Limits: corev1.ResourceList{
                            "cpu":            resources.CPU,
                            "memory":         resources.Memory,
                            "nvidia.com/gpu": resources.GPU,
                        },
                    },
                    VolumeMounts: []corev1.VolumeMount{
                        {
                            Name:      "code-storage",
                            MountPath: "/gemini/code",
                        },
                        {
                            Name:      "output-storage",
                            MountPath: "/gemini/output",
                        },
                    },
                },
            },
            Volumes: []corev1.Volume{
                {
                    Name: "code-storage",
                    VolumeSource: corev1.VolumeSource{
                        PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
                            ClaimName: fmt.Sprintf("code-%s", envID),
                        },
                    },
                },
                {
                    Name: "output-storage",
                    VolumeSource: corev1.VolumeSource{
                        PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
                            ClaimName: fmt.Sprintf("output-%s", envID),
                        },
                    },
                },
            },
        },
    }

    _, err := clientset.CoreV1().Pods("dev-environments").Create(context.TODO(), pod, metav1.CreateOptions{})
    return err
}

// 创建 NodePort Service
func CreateSSHService(envID string, nodePort int) error {
    service := &corev1.Service{
        ObjectMeta: metav1.ObjectMeta{
            Name:      fmt.Sprintf("dev-env-%s-ssh", envID),
            Namespace: "dev-environments",
        },
        Spec: corev1.ServiceSpec{
            Type: corev1.ServiceTypeNodePort,
            Selector: map[string]string{
                "app":    "dev-env",
                "env-id": envID,
            },
            Ports: []corev1.ServicePort{
                {
                    Name:       "ssh",
                    Protocol:   corev1.ProtocolTCP,
                    Port:       22,
                    TargetPort: intstr.FromInt(22),
                    NodePort:   int32(nodePort),
                },
            },
        },
    }

    _, err := clientset.CoreV1().Services("dev-environments").Create(context.TODO(), service, metav1.CreateOptions{})
    return err
}
```

### 9.4 API 实现

```go
package api

// 创建开发环境 API
func CreateEnvironment(c *gin.Context) {
    var req CreateEnvRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 1. 验证用户配额
    customer := GetCurrentCustomer(c)
    if !HasEnoughQuota(customer, req.Resources) {
        c.JSON(403, gin.H{"error": "配额不足"})
        return
    }

    // 2. 分配端口
    port, err := portPool.Allocate()
    if err != nil {
        c.JSON(500, gin.H{"error": "端口分配失败"})
        return
    }

    // 3. 生成 SSH 密码
    password := GeneratePassword()

    // 4. 生成环境 ID
    envID := GenerateEnvID()

    // 5. 创建数据库记录
    env := &Environment{
        ID:         envID,
        CustomerID: customer.ID,
        Name:       req.Name,
        Image:      req.Image,
        Status:     "creating",
        SSHPort:    port,
        SSHPassword: password,
    }
    if err := db.Create(env).Error; err != nil {
        portPool.Release(port)
        c.JSON(500, gin.H{"error": "数据库创建失败"})
        return
    }

    // 6. 创建 Kubernetes 资源
    go func() {
        // 创建 PVC
        if err := k8s.CreatePVC(envID, req.Resources.Storage); err != nil {
            UpdateEnvStatus(envID, "failed")
            return
        }

        // 创建 Pod
        if err := k8s.CreateDevEnvironmentPod(envID, req.Image, password, req.Resources); err != nil {
            UpdateEnvStatus(envID, "failed")
            return
        }

        // 创建 Service
        if err := k8s.CreateSSHService(envID, port); err != nil {
            UpdateEnvStatus(envID, "failed")
            return
        }

        // 等待 Pod Ready
        if err := k8s.WaitForPodReady(envID, 5*time.Minute); err != nil {
            UpdateEnvStatus(envID, "failed")
            return
        }

        // 更新状态为 running
        UpdateEnvStatus(envID, "running")
    }()

    // 7. 返回连接信息
    c.JSON(200, gin.H{
        "env_id":       envID,
        "ssh_host":     config.SSHHost,
        "ssh_port":     port,
        "ssh_username": "developer",
        "ssh_password": password,
        "jupyter_url":  fmt.Sprintf("https://%s/jupyter/%s", config.Domain, envID),
        "status":       "creating",
    })
}

// 获取 SSH 连接信息
func GetSSHInfo(c *gin.Context) {
    envID := c.Param("id")

    var env Environment
    if err := db.Where("id = ? AND customer_id = ?", envID, GetCurrentCustomer(c).ID).First(&env).Error; err != nil {
        c.JSON(404, gin.H{"error": "环境不存在"})
        return
    }

    c.JSON(200, gin.H{
        "ssh_host":     config.SSHHost,
        "ssh_port":     env.SSHPort,
        "ssh_username": "developer",
        "ssh_password": env.SSHPassword,
        "ssh_command":  fmt.Sprintf("ssh developer@%s -p %d", config.SSHHost, env.SSHPort),
    })
}

// 重置 SSH 密码
func ResetSSHPassword(c *gin.Context) {
    envID := c.Param("id")

    var env Environment
    if err := db.Where("id = ? AND customer_id = ?", envID, GetCurrentCustomer(c).ID).First(&env).Error; err != nil {
        c.JSON(404, gin.H{"error": "环境不存在"})
        return
    }

    // 生成新密码
    newPassword := GeneratePassword()

    // 更新数据库
    db.Model(&env).Update("ssh_password", newPassword)

    // 更新容器中的密码
    k8s.ExecInPod(env.ID, fmt.Sprintf("echo 'developer:%s' | chpasswd", newPassword))

    c.JSON(200, gin.H{
        "ssh_password": newPassword,
    })
}
```

---

## 10. 监控和日志

### 10.1 SSH 连接监控

```go
// 记录 SSH 连接
type SSHConnection struct {
    EnvID       string
    CustomerID  int64
    IPAddress   string
    ConnectedAt time.Time
    DisconnectedAt *time.Time
    BytesIn     int64
    BytesOut    int64
}

// 监控脚本（在容器内运行）
#!/bin/bash
# /usr/local/bin/monitor-ssh.sh

while true; do
    # 获取当前 SSH 连接
    who | grep pts | while read line; do
        user=$(echo $line | awk '{print $1}')
        pts=$(echo $line | awk '{print $2}')
        ip=$(echo $line | awk '{print $5}' | tr -d '()')

        # 发送到监控系统
        curl -X POST http://monitor-api/ssh-connections \
            -d "{\"env_id\":\"$ENV_ID\",\"user\":\"$user\",\"ip\":\"$ip\"}"
    done
    sleep 60
done
```

### 10.2 审计日志

```bash
# 配置 SSH 审计日志
# /etc/ssh/sshd_config
SyslogFacility AUTH
LogLevel VERBOSE

# 日志示例
# 登录成功
Jan 26 10:15:23 sshd[1234]: Accepted password for developer from 1.2.3.4 port 52345 ssh2

# 登录失败
Jan 26 10:15:30 sshd[1235]: Failed password for developer from 1.2.3.4 port 52346 ssh2
```

---

## 11. 故障处理

### 11.1 常见问题

| 问题 | 原因 | 解决方案 |
|------|------|----------|
| 连接超时 | 防火墙/端口未开放 | 检查防火墙规则和 NodePort 配置 |
| 密码错误 | 密码未同步 | 重置密码或检查容器状态 |
| 连接被拒绝 | SSH 服务未启动 | 检查容器日志，重启 Pod |
| 端口已占用 | 端口分配冲突 | 检查端口池，释放或分配新端口 |

### 11.2 健康检查

```yaml
# Pod 健康检查配置
livenessProbe:
  tcpSocket:
    port: 22
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  tcpSocket:
    port: 22
  initialDelaySeconds: 5
  periodSeconds: 5
```

---

## 12. 扩展方案

### 12.1 高可用配置

```yaml
# 使用 LoadBalancer 类型（云环境）
apiVersion: v1
kind: Service
metadata:
  name: ssh-gateway
spec:
  type: LoadBalancer
  loadBalancerIP: x.x.x.x
  ports:
  - name: ssh-range
    protocol: TCP
    port: 30000
    targetPort: 30000
```

### 12.2 迁移到 SSH 网关

当环境数量超过 1000 时，可以迁移到 SSH 网关方案：

```
用户 -> SSH Gateway (Port 22) -> 根据用户名路由 -> 对应的 Pod
```

实现工具：
- [sshpiper](https://github.com/tg123/sshpiper)
- [Teleport](https://goteleport.com/)
- 自研 SSH Proxy

---

**文档结束**

这个技术方案提供了从 MVP 到大规模扩展的完整实现路径。