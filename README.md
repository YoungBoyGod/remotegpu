# RemoteGPU - GPU 云平台系统

> 一个功能完善的 GPU 云平台，提供开发环境管理、资源调度、训练推理等完整功能

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/vue-3.x-4FC08D.svg)](https://vuejs.org/)

---

## 📖 项目简介

RemoteGPU 是一个企业级 GPU 云平台系统，旨在为 AI/ML 开发者提供便捷的 GPU 资源管理和使用体验。

### 核心特性

- 🚀 **快速部署** - 一键创建开发环境，支持 Linux 和 Windows
- 💻 **多种访问方式** - SSH、RDP、JupyterLab、Web 终端
- 🎯 **智能调度** - 多种调度策略，优化资源利用率
- 📊 **完善监控** - 实时监控 GPU 使用情况、温度、功耗
- 🔐 **安全可靠** - 多租户隔离、RBAC 权限控制
- 💰 **灵活计费** - 支持按量计费、包年包月等多种模式


---

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                        前端层                                 │
│                   Vue 3 + Element Plus                       │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    API 网关层                                 │
│                  Nginx + Traefik                             │
└────────────────────────┬────────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
        ▼                ▼                ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ 用户管理      │  │ 环境管理      │  │ 资源调度      │
│ 模块         │  │ 模块         │  │ 模块         │
└──────────────┘  └──────────────┘  └──────────────┘
        │                │                │
        └────────────────┼────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    数据层                                     │
│         PostgreSQL + Redis + Etcd + MinIO                   │
└─────────────────────────────────────────────────────────────┘
```


---

## 🎯 核心功能

### 1. 用户管理
- 用户注册与登录（支持 OAuth2、LDAP）
- 用户分类：管理人员、内部用户、外部用户
- 工作空间管理（个人/团队/企业）
- 资源配额管理

### 2. 环境管理
- 快速创建开发环境（< 2 分钟）
- 支持多种镜像（PyTorch、TensorFlow、CUDA）
- SSH/RDP 远程访问
- JupyterLab 集成

### 3. 资源调度
- 智能调度算法（First Fit、Best Fit、Least Loaded）
- GPU 自动分配与释放
- 端口动态管理
- 资源使用监控


### 4. 数据管理
- 数据集上传与版本管理
- 数据集挂载到环境
- 对象存储集成（MinIO/S3）

### 5. 镜像管理
- 官方镜像库
- 自定义镜像构建
- Harbor 镜像仓库集成

### 6. 监控告警
- 实时监控（CPU、内存、GPU、网络）
- 告警规则配置
- Prometheus + Grafana 集成

### 7. 计费管理
- 多种计费模式（按量、包年包月）
- 账单生成与查询
- 支付集成


---

## 🛠️ 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL 14+
- **缓存**: Redis 7+
- **配置中心**: Etcd 3.5+

### 前端
- **框架**: Vue 3
- **UI 库**: Element Plus
- **状态管理**: Pinia
- **构建工具**: Vite


### 容器化
- **容器引擎**: Docker
- **编排平台**: Kubernetes
- **镜像仓库**: Harbor

### 基础设施
- **对象存储**: MinIO / S3
- **配置中心**: Etcd
- **堡垒机**: JumpServer
- **远程桌面**: Guacamole
- **监控**: Prometheus + Grafana


---

## 📋 基础设施要求

### 硬件要求

**最小配置（开发环境）：**
- CPU: 8 核
- 内存: 16GB
- 存储: 100GB SSD
- GPU: 可选

**推荐配置（生产环境）：**
- CPU: 32 核+
- 内存: 128GB+
- 存储: 1TB+ SSD
- GPU: NVIDIA GPU（支持 CUDA）

### 软件依赖

- **操作系统**: Linux (Ubuntu 20.04+ / CentOS 8+)
- **Docker**: 20.10+
- **Kubernetes**: 1.24+
- **PostgreSQL**: 14+
- **Redis**: 7+
- **Etcd**: 3.5+
- **MinIO**: RELEASE.2023-01-01+


---

## 🚀 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/your-org/remotegpu.git
cd remotegpu
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库、Redis 等连接信息
vim .env
```

### 3. 启动基础设施

```bash
# 使用 Docker Compose 启动基础设施
docker-compose up -d postgres redis etcd minio
```

### 4. 初始化数据库

```bash
# 运行数据库迁移
make migrate-up
```

### 5. 启动后端服务

```bash
cd backend
go mod download
go run cmd/server/main.go
```

### 6. 启动前端服务

```bash
cd frontend
npm install
npm run dev
```

### 7. 访问系统

- **前端界面**: http://localhost:3000
- **API 文档**: http://localhost:8080/swagger
- **监控面板**: http://localhost:9090


---

## 📦 部署指南

### Docker 部署

```bash
# 构建镜像
make docker-build

# 启动所有服务
docker-compose up -d
```

### Kubernetes 部署

```bash
# 应用 Kubernetes 配置
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/deployments/
kubectl apply -f k8s/services/
```

详细部署文档请参考：[docs/infrastructure/infrastructure.md](./docs/infrastructure/infrastructure.md)


---

## 📚 文档

### 设计文档
- [系统架构设计](./docs/design/system_architecture.md)
- [模块划分](./docs/design/module_division.md)
- [数据库设计](./docs/design/database_design.md)
- [客户管理设计](./docs/design/customer_management.md)

### 需求文档
- [需求文档索引](./docs/requirements/README.md)
- [总需求文档](./docs/REQUIREMENTS.md)

### 基础设施
- [基础设施配置](./docs/infrastructure/infrastructure.md)
- [第三方系统集成](./docs/infrastructure/third_party_integration.md)


---

## 🔧 开发指南

### 项目结构

```
remotegpu/
├── backend/              # 后端代码
│   ├── cmd/             # 命令行入口
│   ├── internal/        # 内部包
│   │   ├── api/        # API 处理器
│   │   ├── models/     # 数据模型
│   │   ├── services/   # 业务逻辑
│   │   └── repository/ # 数据访问层
│   └── pkg/            # 公共包
├── frontend/            # 前端代码
│   ├── src/
│   │   ├── views/      # 页面组件
│   │   ├── components/ # 通用组件
│   │   ├── api/        # API 调用
│   │   └── stores/     # 状态管理
│   └── public/
├── docs/               # 文档
├── k8s/                # Kubernetes 配置
├── scripts/            # 脚本工具
└── docker-compose.yml  # Docker Compose 配置
```

### 开发环境设置

1. **安装 Go 1.21+**
   ```bash
   # 下载并安装 Go
   wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
   export PATH=$PATH:/usr/local/go/bin
   ```

2. **安装 Node.js 18+**
   ```bash
   # 使用 nvm 安装 Node.js
   curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
   nvm install 18
   nvm use 18
   ```

3. **安装开发工具**
   ```bash
   # 安装 air（Go 热重载）
   go install github.com/cosmtrek/air@latest
   
   # 安装 swag（API 文档生成）
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

### 代码规范

- **Go 代码**: 遵循 [Effective Go](https://golang.org/doc/effective_go.html) 规范
- **Vue 代码**: 遵循 [Vue 3 风格指南](https://vuejs.org/style-guide/)
- **提交信息**: 遵循 [Conventional Commits](https://www.conventionalcommits.org/)

### 运行测试

```bash
# 后端测试
cd backend
go test ./...

# 前端测试
cd frontend
npm run test
```


---

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 贡献流程

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

### 提交规范

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type 类型：**
- `feat`: 新功能
- `fix`: 修复 Bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具链相关

**示例：**
```
feat(environment): 添加环境自动停止功能

- 添加定时任务检查空闲环境
- 超过 2 小时无活动自动停止
- 发送通知给用户

Closes #123
```


---

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件


---

## 📞 联系我们

- **项目主页**: https://github.com/your-org/remotegpu
- **问题反馈**: https://github.com/your-org/remotegpu/issues
- **邮箱**: support@remotegpu.com
- **文档**: https://docs.remotegpu.com


---

## 🙏 致谢

感谢以下开源项目：

- [Gin](https://github.com/gin-gonic/gin) - Go Web 框架
- [Vue.js](https://github.com/vuejs/vue) - 前端框架
- [Element Plus](https://github.com/element-plus/element-plus) - UI 组件库
- [Kubernetes](https://github.com/kubernetes/kubernetes) - 容器编排平台
- [PostgreSQL](https://www.postgresql.org/) - 数据库
- [Redis](https://redis.io/) - 缓存系统


---

**⭐ 如果这个项目对你有帮助，请给我们一个 Star！**
