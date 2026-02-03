# SSH 公网访问部署方案对比

> 传统架构 vs Kubernetes 架构
>
> 创建日期：2026-01-26

---

## 目录

1. [方案总览](#1-方案总览)
2. [方案 1：传统架构（无 K8s）](#2-方案-1传统架构无-k8s)
3. [方案 2：Kubernetes 架构](#3-方案-2kubernetes-架构)
4. [方案对比](#4-方案对比)
5. [选型建议](#5-选型建议)

---

## 1. 方案总览

### 1.1 两种架构对比

```
┌─────────────────────────────────────────────────────────────┐
│                   方案 1：传统架构                             │
├─────────────────────────────────────────────────────────────┤
│  用户 -> 负载均衡 -> Docker 容器 (直接运行在物理机/VM上)        │
│  - 使用 Docker + Docker Compose                             │
│  - 手动管理容器生命周期                                       │
│  - 适合小规模、传统运维团队                                   │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                   方案 2：Kubernetes 架构                      │
├─────────────────────────────────────────────────────────────┤
│  用户 -> 负载均衡 -> K8s Service -> Pod                       │
│  - 使用 Kubernetes 调度                                      │
│  - 自动化管理、自愈、扩展                                     │
│  - 适合大规模、云原生团队                                     │
└─────────────────────────────────────────────────────────────┘
```

---

## 2. 方案 1：传统架构（无 K8s）

### 2.1 架构设计

```
                    Internet
                        │
                        ▼
        ┌───────────────────────────┐
        │    负载均衡器 / 网关        │
        │   (Nginx / HAProxy)       │
        │  ssh.example.com          │
        └────────────┬──────────────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
        ▼            ▼            ▼
    ┌──────┐    ┌──────┐    ┌──────┐
    │ 主机1 │    │ 主机2 │    │ 主机3 │
    │GPU x4│    │GPU x4│    │GPU x8│
    └──┬───┘    └──┬───┘    └──┬───┘
       │           │           │
   ┌───┴───┐   ┌──┴────┐  ┌───┴───┐
   │Docker │   │Docker │  │Docker │
   │Daemon │   │Daemon │  │Daemon │
   └───┬───┘   └──┬────┘  └───┬───┘
       │          │           │
   ┌───┴───┐  ┌──┴────┐  ┌───┴───┐
   │容器池  │  │容器池  │  │容器池  │
   │dev-1  │  │dev-10 │  │dev-20 │
   │dev-2  │  │dev-11 │  │dev-21 │
   │dev-3  │  │dev-12 │  │dev-22 │
   └───────┘  └───────┘  └───────┘
```

### 2.2 核心组件

#### A. 管理服务（Go/Python）

**职责：**
- API 服务器（处理用户请求）
- 容器管理器（创建/删除/监控 Docker 容器）
- 端口分配器（管理 SSH 端口池）
- 调度器（选择主机创建容器）

#### B. 主机管理 Agent

**每台主机运行一个 Agent，负责：**
- 接收管理服务的指令
- 管理本机的 Docker 容器
- 上报主机资源使用情况
- 容器健康检查

#### C. Docker 引擎

**每台主机的 Docker 服务，负责：**
- 运行开发环境容器
- 管理容器网络
- 管理存储卷

#### D. 共享存储（可选）

**网络存储系统：**
- NFS / GlusterFS / Ceph
- 提供持久化存储
- 多主机共享用户数据

### 2.3 实现架构

```
┌─────────────────────────────────────────────────────────────┐
│                     管理平台 (Master)                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  API 服务器   │  │  调度器       │  │  数据库       │     │
│  │  (Gin/Flask) │  │  (Scheduler)  │  │  (Postgres)  │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  端口管理器   │  │  监控服务     │  │  计费服务     │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│                                                             │
└────────────────────────┬────────────────────────────────────┘
                         │ RPC / HTTP
                         │
        ┌────────────────┼────────────────┐
        │                │                │
        ▼                ▼                ▼
┌───────────────┐ ┌───────────────┐ ┌───────────────┐
│   主机 Agent  │ │   主机 Agent  │ │   主机 Agent  │
│   (Worker 1)  │ │   (Worker 2)  │ │   (Worker 3)  │
├───────────────┤ ├───────────────┤ ├───────────────┤
│ Docker Daemon │ │ Docker Daemon │ │ Docker Daemon │
├───────────────┤ ├───────────────┤ ├───────────────┤
│  容器 1-9     │ │  容器 10-19   │ │  容器 20-29   │
└───────────────┘ └───────────────┘ └───────────────┘
```

### 2.4 技术实现

#### 2.4.1 主机 Agent 设计

```go
// worker-agent/main.go
package main

import (
    "github.com/docker/docker/client"
)

type WorkerAgent struct {
    HostID       string
    MasterAddr   string
    DockerClient *client.Client
    Resources    HostResources
}

type HostResources struct {
    TotalCPU    int
    TotalMemory int64
    TotalGPU    int
    UsedCPU     int
    UsedMemory  int64
    UsedGPU     int
}

// 启动 Agent
func (a *WorkerAgent) Start() {
    // 1. 连接到 Docker Daemon
    a.connectDocker()

    // 2. 注册到 Master
    a.registerToMaster()

    // 3. 启动心跳
    go a.heartbeat()

    // 4. 监听 Master 指令
    a.listenCommands()
}

// 创建容器
func (a *WorkerAgent) CreateContainer(req CreateContainerRequest) error {
    // 1. 拉取镜像
    if err := a.pullImage(req.Image); err != nil {
        return err
    }

    // 2. 创建容器
    containerID, err := a.DockerClient.ContainerCreate(
        ctx,
        &container.Config{
            Image: req.Image,
            Env: []string{
                fmt.Sprintf("SSH_PASSWORD=%s", req.SSHPassword),
                "ENABLE_JUPYTER=true",
            },
            ExposedPorts: nat.PortSet{
                "22/tcp":   struct{}{},
                "8888/tcp": struct{}{},
            },
        },
        &container.HostConfig{
            PortBindings: nat.PortMap{
                "22/tcp": []nat.PortBinding{
                    {HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", req.SSHPort)},
                },
            },
            Resources: container.Resources{
                NanoCPUs: req.CPU * 1e9,
                Memory:   req.Memory * 1024 * 1024 * 1024,
                DeviceRequests: []container.DeviceRequest{
                    {
                        Count:        req.GPU,
                        Capabilities: [][]string{{"gpu"}},
                    },
                },
            },
            Binds: []string{
                fmt.Sprintf("%s:/gemini/code", req.CodeVolume),
                fmt.Sprintf("%s:/gemini/output", req.OutputVolume),
            },
        },
        nil,
        nil,
        req.ContainerName,
    )

    if err != nil {
        return err
    }

    // 3. 启动容器
    return a.DockerClient.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

// 删除容器
func (a *WorkerAgent) DeleteContainer(containerID string) error {
    return a.DockerClient.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
        Force: true,
    })
}

// 上报资源使用情况
func (a *WorkerAgent) ReportResources() {
    resources := a.collectResources()

    // 上报到 Master
    http.Post(
        fmt.Sprintf("%s/api/agents/%s/resources", a.MasterAddr, a.HostID),
        "application/json",
        bytes.NewBuffer(resourcesJSON),
    )
}

// 心跳
func (a *WorkerAgent) heartbeat() {
    ticker := time.NewTicker(30 * time.Second)
    for range ticker.C {
        http.Post(
            fmt.Sprintf("%s/api/agents/%s/heartbeat", a.MasterAddr, a.HostID),
            "application/json",
            nil,
        )
    }
}
```

#### 2.4.2 调度器设计

```go
// scheduler/scheduler.go
package scheduler

type Scheduler struct {
    db    *gorm.DB
    hosts map[string]*Host
}

type Host struct {
    ID           string
    IPAddress    string
    Status       string // online, offline, maintenance
    Resources    HostResources
    LastHeartbeat time.Time
}

// 选择主机（简单的资源匹配）
func (s *Scheduler) SelectHost(req ResourceRequirement) (*Host, error) {
    var candidates []*Host

    // 1. 筛选可用主机
    for _, host := range s.hosts {
        if host.Status != "online" {
            continue
        }

        // 检查资源是否满足
        if host.Resources.AvailableCPU >= req.CPU &&
           host.Resources.AvailableMemory >= req.Memory &&
           host.Resources.AvailableGPU >= req.GPU {
            candidates = append(candidates, host)
        }
    }

    if len(candidates) == 0 {
        return nil, errors.New("no available host")
    }

    // 2. 选择策略（可以实现多种策略）
    // 策略 1：最少使用（Least Used）
    sort.Slice(candidates, func(i, j int) bool {
        return candidates[i].Resources.UsedCPU < candidates[j].Resources.UsedCPU
    })

    return candidates[0], nil
}

// 创建开发环境
func (s *Scheduler) CreateEnvironment(req CreateEnvRequest) (*Environment, error) {
    // 1. 选择主机
    host, err := s.SelectHost(req.Resources)
    if err != nil {
        return nil, err
    }

    // 2. 分配端口
    port, err := s.AllocatePort()
    if err != nil {
        return nil, err
    }

    // 3. 生成配置
    envID := GenerateEnvID()
    password := GeneratePassword()

    // 4. 调用 Agent 创建容器
    agentAddr := fmt.Sprintf("http://%s:8080", host.IPAddress)
    resp, err := http.Post(
        fmt.Sprintf("%s/containers/create", agentAddr),
        "application/json",
        createContainerRequestJSON,
    )

    if err != nil {
        return nil, err
    }

    // 5. 保存到数据库
    env := &Environment{
        ID:          envID,
        CustomerID:  req.CustomerID,
        HostID:      host.ID,
        SSHPort:     port,
        SSHPassword: password,
        Status:      "running",
    }
    s.db.Create(env)

    return env, nil
}
```

#### 2.4.3 端口映射实现

**方式 1：直接端口映射（简单）**

```bash
# Docker 直接映射主机端口
docker run -d \
  -p 30001:22 \
  -p 38001:8888 \
  --name dev-env-123 \
  dev-image:latest

# 用户访问
ssh developer@主机1IP:30001
ssh developer@主机2IP:30002
```

**方式 2：负载均衡器转发（小规模或配合自动化）**

```nginx
# /etc/nginx/stream.conf
stream {
    # SSH 端口映射
    upstream ssh_30001 {
        server 192.168.1.10:30001;  # 主机1
    }

    upstream ssh_30002 {
        server 192.168.1.10:30002;  # 主机1
    }

    upstream ssh_30003 {
        server 192.168.1.11:30001;  # 主机2
    }

    server {
        listen 30001;
        proxy_pass ssh_30001;
    }

    server {
        listen 30002;
        proxy_pass ssh_30002;
    }

    server {
        listen 30003;
        proxy_pass ssh_30003;
    }
}
```

> 端口规模大时需要配合自动化生成配置并热更新（如模板渲染 + reload、HAProxy Runtime API 等），否则不建议作为默认方案。

**方式 3：动态端口转发（最灵活）**

```go
// 使用 iptables 动态配置端口转发
func ConfigurePortForwarding(externalPort, hostIP string, hostPort int) error {
    // 添加 DNAT 规则
    cmd := fmt.Sprintf(
        "iptables -t nat -A PREROUTING -p tcp --dport %d -j DNAT --to-destination %s:%d",
        externalPort, hostIP, hostPort,
    )
    return exec.Command("sh", "-c", cmd).Run()
}

// 删除端口转发
func RemovePortForwarding(externalPort int, hostIP string, hostPort int) error {
    cmd := fmt.Sprintf(
        "iptables -t nat -D PREROUTING -p tcp --dport %d -j DNAT --to-destination %s:%d",
        externalPort, hostIP, hostPort,
    )
    return exec.Command("sh", "-c", cmd).Run()
}
```

#### 2.4.4 存储方案

**方案 A：本地存储**

```bash
# 每台主机本地目录
/data/environments/
  ├── env-123/
  │   ├── code/      # 代码目录
  │   └── output/    # 输出目录
  ├── env-124/
  └── env-125/

# 容器挂载
docker run -v /data/environments/env-123/code:/gemini/code \
           -v /data/environments/env-123/output:/gemini/output \
           ...
```

**优点：** 简单、性能好
**缺点：** 容器只能在当前主机运行，无法迁移

**方案 B：NFS 共享存储**

```bash
# NFS 服务器
/export/environments/
  ├── env-123/
  ├── env-124/
  └── env-125/

# 主机挂载 NFS
mount -t nfs nfs-server:/export/environments /mnt/environments

# 容器挂载
docker run -v /mnt/environments/env-123/code:/gemini/code \
           ...
```

**优点：** 容器可以在任意主机运行
**缺点：** 网络存储性能较差

**推荐：** MVP 使用方案 A，后期迁移到方案 B

#### 2.4.5 数据库设计

```sql
-- 主机表
CREATE TABLE hosts (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(128),
    ip_address VARCHAR(64) NOT NULL,
    status VARCHAR(20) DEFAULT 'online',  -- online, offline, maintenance
    total_cpu INT,
    total_memory BIGINT,
    total_gpu INT,
    used_cpu INT DEFAULT 0,
    used_memory BIGINT DEFAULT 0,
    used_gpu INT DEFAULT 0,
    last_heartbeat TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 环境表
CREATE TABLE environments (
    id VARCHAR(64) PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    host_id VARCHAR(64) NOT NULL,
    container_id VARCHAR(128),
    name VARCHAR(128),
    image VARCHAR(256),
    status VARCHAR(20),  -- creating, running, stopped, failed
    ssh_port INT,
    ssh_password VARCHAR(128),
    jupyter_port INT,
    cpu INT,
    memory BIGINT,
    gpu INT,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (host_id) REFERENCES hosts(id),
    INDEX idx_customer_id (customer_id),
    INDEX idx_host_id (host_id)
);

-- 端口映射表（与 K8s 方案相同）
CREATE TABLE port_mappings (
    id BIGSERIAL PRIMARY KEY,
    env_id VARCHAR(64) NOT NULL,
    external_port INT NOT NULL UNIQUE,
    host_id VARCHAR(64) NOT NULL,
    host_port INT NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    allocated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (env_id) REFERENCES environments(id),
    FOREIGN KEY (host_id) REFERENCES hosts(id)
);
```

### 2.5 部署架构

```yaml
# docker-compose.yml (管理节点)
version: '3.8'

services:
  # API 服务
  api-server:
    build: ./api-server
    ports:
      - "8000:8000"
    environment:
      - DATABASE_URL=postgresql://user:pass@postgres:5432/remotegpu
      - REDIS_URL=redis://redis:6379
    depends_on:
      - postgres
      - redis

  # 调度器
  scheduler:
    build: ./scheduler
    environment:
      - DATABASE_URL=postgresql://user:pass@postgres:5432/remotegpu

  # 数据库
  postgres:
    image: postgres:14
    volumes:
      - postgres-data:/var/lib/postgresql/data
    environment:
      - POSTGRES_DB=remotegpu
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass

  # Redis
  redis:
    image: redis:7
    volumes:
      - redis-data:/data

volumes:
  postgres-data:
  redis-data:
```

```bash
# 工作节点部署脚本
#!/bin/bash
# deploy-worker.sh

# 1. 安装 Docker
curl -fsSL https://get.docker.com | sh

# 2. 安装 NVIDIA Docker Runtime
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | apt-key add -
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | \
    tee /etc/apt/sources.list.d/nvidia-docker.list
apt-get update && apt-get install -y nvidia-docker2
systemctl restart docker

# 3. 下载并启动 Worker Agent
wget https://releases.example.com/worker-agent
chmod +x worker-agent
./worker-agent --master=http://master-ip:8000 --host-id=worker-1
```

### 2.6 完整流程示例

```
1. 用户请求创建环境
   ↓
2. API 服务器接收请求
   ↓
3. 调度器选择合适的主机（主机1）
   ↓
4. 分配端口（30001）
   ↓
5. 调用主机1的 Agent 创建容器
   ↓
6. Agent 使用 Docker API 创建并启动容器
   - 映射端口 30001:22
   - 挂载存储卷
   - 设置 GPU
   ↓
7. 配置 Nginx/iptables 端口转发（如果需要）
   外部 30001 -> 主机1:30001
   ↓
8. 返回连接信息给用户
   ssh developer@ssh.example.com -p 30001
```

---

## 3. 方案 2：Kubernetes 架构

### 3.1 架构设计

```
                    Internet
                        │
                        ▼
        ┌───────────────────────────┐
        │    负载均衡器 / Ingress    │
        │   ssh.example.com          │
        └────────────┬──────────────┘
                     │
        ┌────────────┼────────────┐
        │   Kubernetes 集群        │
        └────────────┬──────────────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
        ▼            ▼            ▼
    ┌──────┐    ┌──────┐    ┌──────┐
    │ Node1│    │ Node2│    │ Node3│
    │GPU x4│    │GPU x4│    │GPU x8│
    └──┬───┘    └──┬───┘    └──┬───┘
       │           │           │
   ┌───┴───┐   ┌──┴────┐  ┌───┴───┐
   │Pod 1  │   │Pod 10 │  │Pod 20 │
   │Pod 2  │   │Pod 11 │  │Pod 21 │
   │Pod 3  │   │Pod 12 │  │Pod 22 │
   └───────┘   └───────┘  └───────┘
```

### 3.2 核心组件

#### A. Kubernetes 原生组件

- **API Server**: 统一入口
- **Scheduler**: 自动调度 Pod
- **Controller Manager**: 管理 Pod 生命周期
- **kubelet**: 每个节点的 Agent
- **GPU Operator**: GPU 资源管理

#### B. 自研组件

- **API 服务**: 业务逻辑和用户接口
- **Operator**: 自定义资源管理（可选）
- **监控服务**: 资源监控和计费

### 3.3 技术实现

#### 3.3.1 Pod 创建

```yaml
# dev-environment-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: dev-env-123
  namespace: dev-environments
  labels:
    app: dev-env
    env-id: env-123
    customer-id: "456"
spec:
  containers:
  - name: dev-env
    image: registry.example.com/dev-ubuntu20-pytorch:latest
    env:
    - name: SSH_PASSWORD
      valueFrom:
        secretKeyRef:
          name: ssh-secret-123
          key: password
    - name: ENABLE_JUPYTER
      value: "true"
    ports:
    - name: ssh
      containerPort: 22
    - name: jupyter
      containerPort: 8888
    resources:
      limits:
        cpu: "4"
        memory: "16Gi"
        nvidia.com/gpu: "1"
      requests:
        cpu: "1"
        memory: "4Gi"
        nvidia.com/gpu: "1"
    volumeMounts:
    - name: code-storage
      mountPath: /gemini/code
    - name: output-storage
      mountPath: /gemini/output
  volumes:
  - name: code-storage
    persistentVolumeClaim:
      claimName: code-pvc-123
  - name: output-storage
    persistentVolumeClaim:
      claimName: output-pvc-123
```

#### 3.3.2 Service 创建

```yaml
# dev-environment-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: dev-env-123-ssh
  namespace: dev-environments
spec:
  type: NodePort
  selector:
    app: dev-env
    env-id: env-123
  ports:
  - name: ssh
    protocol: TCP
    port: 22
    targetPort: 22
    nodePort: 30001
```

#### 3.3.3 Go 代码实现

```go
// k8s/client.go
package k8s

import (
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

type K8sManager struct {
    clientset *kubernetes.Clientset
}

// 创建开发环境
func (m *K8sManager) CreateEnvironment(req CreateEnvRequest) (*Environment, error) {
    // 1. 创建 Secret（存储 SSH 密码）
    secret := &corev1.Secret{
        ObjectMeta: metav1.ObjectMeta{
            Name:      fmt.Sprintf("ssh-secret-%s", req.EnvID),
            Namespace: "dev-environments",
        },
        StringData: map[string]string{
            "password": req.SSHPassword,
        },
    }
    _, err := m.clientset.CoreV1().Secrets("dev-environments").Create(ctx, secret, metav1.CreateOptions{})
    if err != nil {
        return nil, err
    }

    // 2. 创建 PVC
    if err := m.createPVC(req.EnvID, req.Storage); err != nil {
        return nil, err
    }

    // 3. 创建 Pod
    pod := m.buildPodSpec(req)
    _, err = m.clientset.CoreV1().Pods("dev-environments").Create(ctx, pod, metav1.CreateOptions{})
    if err != nil {
        return nil, err
    }

    // 4. 创建 Service
    service := m.buildServiceSpec(req.EnvID, req.SSHPort)
    _, err = m.clientset.CoreV1().Services("dev-environments").Create(ctx, service, metav1.CreateOptions{})
    if err != nil {
        return nil, err
    }

    // 5. 等待 Pod Ready
    if err := m.waitForPodReady(req.EnvID, 5*time.Minute); err != nil {
        return nil, err
    }

    return &Environment{
        ID:     req.EnvID,
        Status: "running",
        SSHPort: req.SSHPort,
    }, nil
}

// 删除环境
func (m *K8sManager) DeleteEnvironment(envID string) error {
    namespace := "dev-environments"

    // 删除 Service
    m.clientset.CoreV1().Services(namespace).Delete(
        ctx,
        fmt.Sprintf("dev-env-%s-ssh", envID),
        metav1.DeleteOptions{},
    )

    // 删除 Pod
    m.clientset.CoreV1().Pods(namespace).Delete(
        ctx,
        fmt.Sprintf("dev-env-%s", envID),
        metav1.DeleteOptions{},
    )

    // 删除 Secret
    m.clientset.CoreV1().Secrets(namespace).Delete(
        ctx,
        fmt.Sprintf("ssh-secret-%s", envID),
        metav1.DeleteOptions{},
    )

    return nil
}
```

### 3.4 部署架构

```yaml
# helm/values.yaml
# 使用 Helm 部署应用

api:
  replicaCount: 3
  image:
    repository: registry.example.com/remotegpu-api
    tag: v1.0.0
  resources:
    limits:
      cpu: 2
      memory: 4Gi

postgres:
  enabled: true
  persistence:
    size: 100Gi

redis:
  enabled: true
  cluster:
    enabled: true
    nodes: 3
```

```bash
# 部署命令
helm install remotegpu ./helm \
  --namespace remotegpu-system \
  --create-namespace
```

---

## 4. 方案对比

### 4.1 功能对比

| 功能 | 传统架构 | Kubernetes 架构 |
|------|---------|----------------|
| **容器编排** | 手动 Docker 管理 | Kubernetes 自动调度 |
| **资源调度** | 自研调度器 | K8s Scheduler |
| **自动扩缩容** | 手动扩容 | HPA 自动扩容 |
| **自愈能力** | 需要监控脚本 | K8s 自动重启 |
| **负载均衡** | Nginx/HAProxy | K8s Service/Ingress |
| **服务发现** | 手动配置 | K8s DNS |
| **配置管理** | 文件/数据库 | ConfigMap/Secret |
| **存储管理** | NFS/本地存储 | PV/PVC 抽象 |
| **网络隔离** | Docker 网络 | Network Policy |
| **滚动更新** | 手动更新 | Deployment 滚动更新 |

### 4.2 运维对比

| 维度 | 传统架构 | Kubernetes 架构 |
|------|---------|----------------|
| **部署复杂度** | ⭐⭐ 简单 | ⭐⭐⭐⭐ 复杂 |
| **运维成本** | ⭐⭐⭐ 中等 | ⭐⭐⭐⭐ 较高 |
| **学习曲线** | ⭐⭐ 低 | ⭐⭐⭐⭐⭐ 高 |
| **可扩展性** | ⭐⭐⭐ 中等 | ⭐⭐⭐⭐⭐ 优秀 |
| **稳定性** | ⭐⭐⭐ 中等 | ⭐⭐⭐⭐⭐ 优秀 |
| **故障恢复** | ⭐⭐ 需要人工介入 | ⭐⭐⭐⭐⭐ 自动恢复 |
| **监控能力** | ⭐⭐⭐ 需要自建 | ⭐⭐⭐⭐ 生态完善 |

### 4.3 成本对比

| 成本项 | 传统架构 | Kubernetes 架构 |
|--------|---------|----------------|
| **硬件成本** | 较低（按需添加服务器） | 较高（需要一定规模才合算） |
| **人力成本** | 中等（需要运维人员） | 较高（需要 K8s 专家） |
| **开发成本** | 中等（需要自研调度器） | 较低（使用 K8s API） |
| **运维成本** | 中等（手动运维） | 较低（自动化程度高） |
| **学习成本** | 低（Docker 知识） | 高（K8s 生态复杂） |

### 4.4 适用场景

#### 传统架构适合：

✅ **小型团队**（5-20 人）
✅ **小规模部署**（< 100 并发环境）
✅ **传统运维团队**（熟悉 Docker，不熟悉 K8s）
✅ **预算有限**（无法投入 K8s 学习和部署成本）
✅ **快速 MVP**（快速验证业务模式）
✅ **自有机房**（物理机/虚拟机部署）

#### Kubernetes 架构适合：

✅ **中大型团队**（> 20 人）
✅ **大规模部署**（> 100 并发环境）
✅ **云原生团队**（熟悉 K8s 生态）
✅ **云上部署**（AWS/阿里云/腾讯云）
✅ **高可用需求**（99.9% 以上 SLA）
✅ **复杂业务**（多种服务、微服务架构）

---

## 5. 选型建议

### 5.1 阶段性演进路线

```
┌────────────────────────────────────────────────────────────┐
│  阶段 1：MVP（3-6 个月）                                      │
│  - 使用传统架构                                              │
│  - 快速验证业务                                              │
│  - 支持 10-50 并发环境                                       │
└────────────────┬───────────────────────────────────────────┘
                 │
                 ▼
┌────────────────────────────────────────────────────────────┐
│  阶段 2：增长期（6-12 个月）                                  │
│  - 继续使用传统架构                                          │
│  - 优化调度算法                                              │
│  - 支持 50-200 并发环境                                      │
│  - 开始规划 K8s 迁移                                         │
└────────────────┬───────────────────────────────────────────┘
                 │
                 ▼
┌────────────────────────────────────────────────────────────┐
│  阶段 3：规模化（12+ 个月）                                   │
│  - 迁移到 Kubernetes                                        │
│  - 支持 200+ 并发环境                                        │
│  - 自动化运维                                                │
│  - 多云/混合云部署                                           │
└────────────────────────────────────────────────────────────┘
```

### 5.2 决策树

```
开始
  │
  ├─ 团队是否熟悉 K8s？
  │   ├─ 是 → 是否有 K8s 集群？
  │   │         ├─ 是 → 选择 K8s 架构 ✅
  │   │         └─ 否 → 评估部署成本
  │   │                  ├─ 可接受 → K8s 架构
  │   │                  └─ 不可接受 → 传统架构
  │   │
  │   └─ 否 → 是否有学习和投入意愿？
  │            ├─ 是 → 预计用户规模？
  │            │       ├─ > 100 → K8s 架构（值得投入）
  │            │       └─ < 100 → 传统架构（先快速启动）
  │            │
  │            └─ 否 → 传统架构 ✅
  │
  └─ 是否在云上部署？
      ├─ 是 → 选择 K8s 架构（云厂商托管 K8s）
      └─ 否 → 评估规模
              ├─ 小规模 → 传统架构
              └─ 大规模 → K8s 架构
```

### 5.3 迁移路径

如果选择从传统架构迁移到 K8s：

```
阶段 1：准备（1-2 月）
├─ 搭建 K8s 测试集群
├─ 团队学习 K8s
├─ 改造应用适配 K8s
└─ 灰度测试

阶段 2：迁移（2-3 月）
├─ 新用户优先使用 K8s
├─ 逐步迁移老用户
├─ 双架构并行运行
└─ 数据迁移

阶段 3：完成（1 月）
├─ 关闭传统架构
├─ 全部迁移到 K8s
└─ 优化和监控
```

---

## 6. 实施建议

### 6.1 传统架构实施步骤

**第 1 周：基础设施**
- [ ] 采购/准备服务器（至少 2 台 GPU 服务器 + 1 台管理服务器）
- [ ] 安装操作系统（Ubuntu 20.04/22.04）
- [ ] 安装 Docker + NVIDIA Docker Runtime
- [ ] 配置网络和防火墙

**第 2 周：核心服务**
- [ ] 开发 API 服务器（Go/Python）
- [ ] 开发调度器
- [ ] 开发 Worker Agent
- [ ] 搭建数据库（PostgreSQL + Redis）

**第 3 周：容器化**
- [ ] 制作开发环境 Docker 镜像（SSH + JupyterLab）
- [ ] 测试容器创建/删除流程
- [ ] 实现端口映射

**第 4 周：测试和优化**
- [ ] 功能测试
- [ ] 性能测试
- [ ] 监控和日志
- [ ] 上线

### 6.2 Kubernetes 架构实施步骤

**第 1 个月：集群搭建**
- [ ] 规划集群架构（节点数量、网络方案）
- [ ] 安装 Kubernetes 集群
- [ ] 安装 GPU Operator
- [ ] 安装监控系统（Prometheus + Grafana）
- [ ] 配置存储（Ceph/NFS/云存储）

**第 2 个月：应用开发**
- [ ] 开发 API 服务（对接 K8s API）
- [ ] 制作容器镜像
- [ ] 编写 Helm Chart
- [ ] 测试 Pod/Service 创建

**第 3 个月：测试和优化**
- [ ] 功能测试
- [ ] 性能测试
- [ ] 灰度发布
- [ ] 正式上线

---

## 7. 总结

### 核心建议

1. **初创/小团队** → 传统架构（Docker + 自研调度器）
   - 快速启动，低学习成本
   - 3-6 个月可以上线 MVP

2. **中型团队/云上部署** → Kubernetes 架构
   - 长期投入回报高
   - 可扩展性和稳定性更好

3. **大型企业/已有 K8s** → 直接 Kubernetes 架构
   - 充分利用现有基础设施
   - 统一的运维体系

### 关键成功因素

无论选择哪种方案：
- ✅ 端口管理要做好（避免冲突）
- ✅ 安全机制要完善（密码策略、访问控制）
- ✅ 监控要到位（资源使用、SSH 连接）
- ✅ 容器镜像要优化（启动速度、大小）
- ✅ 存储方案要合理（性能 vs 灵活性）

---

**文档结束**

根据你的实际情况选择合适的架构，并做好长期演进规划。
