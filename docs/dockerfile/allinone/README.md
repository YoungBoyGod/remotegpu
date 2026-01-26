# GPU Workspace 容器管理工具

一套完整的 GPU 工作空间容器化解决方案，支持 SSH、Jupyter Lab 和 VSCode Web。

## 📁 目录结构

```
allinone/
├── Dockerfile              # 容器镜像定义
├── docker-compose.yml      # 容器编排配置（已配置资源限制）
├── entrypoint.sh          # 容器启动脚本
├── get_ssh_key.sh         # SSH 密钥获取工具 ⭐
├── configure_limits.sh    # 资源限制配置工具 ⭐
├── data/                  # 用户数据目录（挂载卷）
│   └── user001/
└── ssh_keys/              # 生成的 SSH 密钥包（运行脚本后生成）
```

## 🚀 快速开始

### 1. 构建并启动容器

```bash
# 构建镜像
docker build -t gpu-workspace:latest .

# 启动容器
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### 2. 获取 SSH 密钥

```bash
# 运行自动化脚本
./get_ssh_key.sh
```

生成的文件在 `./ssh_keys/` 目录：
- `user001_ssh_package.tar.gz` - 完整分发包（发给用户）
- `user001_使用说明.txt` - 详细使用指南

### 3. 配置资源限制（可选）

```bash
# 运行资源配置工具
./configure_limits.sh
```

选择预设场景：
- 小型工作空间（1-2 用户）：8GB 内存 + 4 核 CPU
- 中型工作空间（3-5 用户）：16GB 内存 + 8 核 CPU
- 大型工作空间（5-10 用户）：32GB 内存 + 16 核 CPU
- 生产环境（10+ 用户）：64GB 内存 + 32 核 CPU

## 🎯 功能特性

### ✅ 服务

| 服务 | 端口 | 说明 |
|------|------|------|
| SSH | 2222 | 命令行访问、VSCode Remote |
| Jupyter Lab | 18888 | Web 版 Python 开发环境 |
| VSCode Web | 18080 | 浏览器版 VSCode |
| 端口转发池 | 19000-19010 | SSH 隧道端口转发 |

### ✅ 自动化功能

- 🔐 **自动生成 SSH 密钥** - 首次启动自动创建
- 🔧 **自动修复权限** - 启动时自动修复挂载目录权限
- 📦 **一键打包分发** - 自动生成用户使用包
- ⚙️ **资源限制配置** - 交互式配置内存/CPU/磁盘限制

### ✅ 资源管理

当前配置（可通过 `configure_limits.sh` 修改）：
- 内存：16GB（保证 8GB）
- CPU：8 核
- 磁盘：100GB（可写层，需启用）
- GPU：1 个（可配置）
- 进程：最多 2000 个

## 📚 文档

| 文档 | 说明 |
|------|------|
| `QUICKSTART.md` | 快速开始指南 |
| `SSH_GUIDE.md` | SSH 登录详细说明 |
| `NETWORK_CONFIG.md` | 网络配置方案 |
| `RESOURCE_LIMITS.md` | 资源限制完整文档 |
| `STORAGE_MANAGEMENT.md` | 存储管理最佳实践 ⭐ |
| `STORAGE_EXAMPLES.md` | 存储配置示例 |
| `DEEP_LEARNING_SETUP.md` | 深度学习环境配置 🔥 |

## 🛠️ 管理工具

### get_ssh_key.sh - SSH 密钥获取工具

自动提取和打包 SSH 私钥，生成用户使用说明。

```bash
./get_ssh_key.sh
```

**生成内容：**
- 私钥文件
- 公钥文件
- 详细使用说明（中文）
- SSH 配置模板
- 完整分发压缩包

### configure_limits.sh - 资源限制配置工具

交互式配置容器资源限制。

```bash
./configure_limits.sh
```

**支持配置：**
- 内存限制和保留
- CPU 核心数
- 磁盘空间限制
- GPU 分配
- 进程数限制

### migrate_docker_storage.sh - Docker 存储迁移工具 ⭐

将 Docker 根目录迁移到更大的分区，解决空间不足问题。

```bash
sudo ./migrate_docker_storage.sh
```

**功能特性：**
- 自动检测当前配置
- 空间验证和安全检查
- 数据同步带进度显示
- 自动更新 Docker 配置
- 验证迁移结果

### monitor_storage.sh - 存储监控和清理工具 ⭐

监控 Docker 和宿主机存储使用情况，提供交互式清理。

```bash
./monitor_storage.sh
```

**监控内容：**
- Docker 镜像、容器、卷使用情况
- 宿主机磁盘使用率
- 用户数据目录占用
- 可回收空间分析
- 存储健康检查

**清理选项：**
- 安全清理（推荐）
- 深度清理
- 自定义清理

### rebuild.sh - 镜像重建工具 🔄

更新 Dockerfile 后重建镜像并重启容器。

```bash
./rebuild.sh
```

**使用场景：**
- 修改了 Dockerfile
- 需要添加新的系统依赖
- 更新基础镜像

**功能：**
- 停止现有容器
- 可选删除旧镜像
- 构建新镜像
- 启动容器
- 验证运行状态

## 🔧 常用命令

### 容器管理

```bash
# 启动容器
docker-compose up -d

# 停止容器
docker-compose down

# 重启容器
docker-compose restart

# 查看日志
docker-compose logs -f

# 查看状态
docker-compose ps

# 进入容器
docker exec -it user001-workspace bash
```

### 资源监控

```bash
# 实时监控资源使用
docker stats user001-workspace

# 查看容器详细信息
docker inspect user001-workspace

# 在容器内查看资源
docker exec user001-workspace free -h      # 内存
docker exec user001-workspace df -h        # 磁盘
docker exec user001-workspace nvidia-smi   # GPU
```

### 动态调整资源

```bash
# 调整内存限制
docker update --memory 32g user001-workspace

# 调整 CPU 限制
docker update --cpus 16 user001-workspace

# 调整多个资源
docker update --memory 32g --cpus 16 user001-workspace
```

### 存储管理

```bash
# 查看存储使用
docker system df
docker system df -v

# 运行监控和清理工具
./monitor_storage.sh

# 迁移 Docker 存储（需要 sudo）
sudo ./migrate_docker_storage.sh

# 清理未使用的资源
docker system prune -f              # 安全清理
docker image prune -a -f            # 清理所有未使用镜像
docker volume prune -f              # 清理未使用的卷

# 查看目录占用
du -sh /var/lib/docker/*
du -sh ./data/user001/
```

## 🌐 访问服务

### SSH 连接

```bash
# 使用密钥连接
ssh -i ~/.ssh/workspace_key -p 2222 gpuuser@服务器IP

# VSCode Remote SSH
# 配置 ~/.ssh/config 后直接连接
code --remote ssh-remote+workspace-user001
```

### Web 服务

- Jupyter Lab: http://服务器IP:18888
- VSCode Web: http://服务器IP:18080

## 🔒 安全建议

### 密钥管理

- ⚠️ 通过加密渠道传输私钥
- ⚠️ 不要上传私钥到 Git 仓库
- ⚠️ 定期轮换密钥

### 容器安全

- ✅ 已禁用密码登录（仅密钥认证）
- ✅ 已配置资源限制
- ✅ 使用非特权用户运行服务

### 密钥轮换

```bash
# 删除旧密钥
rm -rf ./data/user001/.ssh/

# 重启容器（自动生成新密钥）
docker-compose restart

# 重新获取密钥
./get_ssh_key.sh
```

## 📊 多用户部署

### 方式 1：复制目录

```bash
# 复制整个目录
cp -r allinone user002-workspace
cd user002-workspace

# 修改配置
# 1. 修改 docker-compose.yml 中的容器名和端口
# 2. 修改 get_ssh_key.sh 中的 USER_ID

# 启动
docker-compose up -d
```

### 方式 2：使用单独的 compose 文件

为每个用户创建独立的 `docker-compose-user002.yml`，使用不同的端口和数据目录。

```bash
docker-compose -f docker-compose-user002.yml up -d
```

## 🐛 故障排查

### 容器无法启动

```bash
# 查看日志
docker-compose logs

# 检查权限
ls -la ./data/user001

# 修复权限
sudo chown -R 1000:1000 ./data/user001
```

### SSH 连接失败

```bash
# 检查私钥权限
chmod 600 ~/.ssh/workspace_key

# 检查容器状态
docker-compose ps

# 查看 SSH 服务
docker exec user001-workspace ps aux | grep sshd
```

### 资源限制不生效

```bash
# 查看当前限制
docker inspect user001-workspace | grep -E "Memory|Cpu"

# 重启容器应用配置
docker-compose down && docker-compose up -d
```

### 磁盘空间限制不生效

```bash
# 检查存储驱动
docker info | grep "Storage Driver"

# 需要是 overlay2 并配置了配额支持
# 参考 RESOURCE_LIMITS.md 文档配置
```

## 🔄 更新和维护

### 更新镜像

```bash
# 重新构建镜像
docker build -t gpu-workspace:latest .

# 重新创建容器
docker-compose up -d --force-recreate
```

### 备份数据

```bash
# 备份用户数据
tar -czf user001_backup_$(date +%Y%m%d).tar.gz ./data/user001/

# 备份配置
cp docker-compose.yml docker-compose.yml.backup
```

### 清理资源

```bash
# 停止并删除容器
docker-compose down

# 删除镜像
docker rmi gpu-workspace:latest

# 清理未使用的资源
docker system prune -a
```

## 📝 配置示例

### 小型单用户工作空间

```yaml
mem_limit: 8g
cpus: 4
storage_opt:
  size: '50G'
environment:
  - NVIDIA_VISIBLE_DEVICES=0
```

### 大型多用户生产环境

```yaml
mem_limit: 64g
cpus: 32
storage_opt:
  size: '500G'
environment:
  - NVIDIA_VISIBLE_DEVICES=all
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

---

## ⚡ 快速参考

```bash
# 一键部署
docker build -t gpu-workspace:latest . && docker-compose up -d

# 获取密钥并配置资源
./get_ssh_key.sh && ./configure_limits.sh

# 监控运行状态
docker stats user001-workspace

# 进入容器
docker exec -it user001-workspace bash
```

🎉 现在你的 GPU 工作空间已经准备就绪！
