# RemoteGPU Backend

基于 Gin 框架和 GoFrame 分层设计的后端项目。

## 项目结构

```
backend/
├── api/v1/              # API 接口定义（请求/响应结构体）
├── cmd/                 # 程序入口
├── config/              # 配置文件
├── internal/            # 内部代码
│   ├── controller/      # 控制器层
│   ├── service/         # 业务逻辑层
│   ├── dao/             # 数据访问层
│   ├── model/           # 数据模型
│   ├── middleware/      # 中间件
│   └── router/          # 路由配置
└── pkg/                 # 公共库
    ├── auth/            # JWT 认证
    ├── database/        # 数据库
    ├── logger/          # 日志
    ├── redis/           # Redis
    ├── response/        # 统一响应
    └── errors/          # 错误处理
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置文件

修改 `config/config.yaml` 中的数据库和 Redis 配置。

### 3. 运行项目

```bash
# 开发模式（热更新）
make dev

# 或直接运行
go run cmd/main.go

# 指定配置文件
go run cmd/main.go --config=./config/config.yaml --mode=debug
```

### 4. 编译

```bash
make build
```

## API 接口

### 用户相关

- `POST /api/v1/user/register` - 用户注册
- `POST /api/v1/user/login` - 用户登录
- `GET /api/v1/user/:id` - 获取用户信息
- `GET /api/v1/user/info` - 获取当前用户信息（需要认证）
- `PUT /api/v1/user/info` - 更新用户信息（需要认证）

### 健康检查

- `GET /api/v1/health` - 健康检查

## 技术栈

- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **缓存**: Redis
- **日志**: Zap
- **认证**: JWT
- **热更新**: Air

## 开发命令

```bash
make help    # 查看所有命令
make build   # 编译项目
make run     # 运行项目
make dev     # 开发模式（热更新）
make clean   # 清理编译文件
make test    # 运行测试
```
